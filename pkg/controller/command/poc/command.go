package poc

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/hyperledger/aries-framework-go/pkg/common/log"
	//"github.com/hyperledger/aries-framework-go/pkg/controller/command/poc/files"
	"github.com/hyperledger/aries-framework-go/pkg/controller/command"
	vcwalletc "github.com/hyperledger/aries-framework-go/pkg/controller/command/vcwallet"
	vdrc "github.com/hyperledger/aries-framework-go/pkg/controller/command/vdr"
	"github.com/hyperledger/aries-framework-go/pkg/controller/internal/cmdutil"
	"github.com/hyperledger/aries-framework-go/pkg/crypto"
	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"github.com/hyperledger/aries-framework-go/pkg/internal/logutil"
	"github.com/hyperledger/aries-framework-go/pkg/vdr"
	"github.com/hyperledger/aries-framework-go/pkg/wallet"
	"github.com/hyperledger/aries-framework-go/spi/storage"
	"github.com/piprate/json-gold/ld"
)

var (
	logger = log.New("aries-framework/command/poc")
)

// Error codes.
const (
	// InvalidRequestErrorCode is typically a code for invalid requests.
	InvalidRequestErrorCode = command.Code(iota + command.POC)
	NewDIDRequestErrorCode
	DoDeviceEnrolmentRequestErrorCode
	GenerateVPRequestErrorCode
	VerifyCredentialRequestErrorCode
	AcceptEnrolmentRequestErrorCode
	TestingCallRequestErrorCode
	GetTrustedIssuerListrRequestErrorCode
	SignJWTContentErrorCode
	StoreCredentialRequestErrorCode
	DeriveProofRequestErrorCode
)

// constants for the VDR controller's methods.
const (
	// command name.
	CommandName = "poc"

	// command methods.
	NewDIDCommandMethod            = "NewDID"
	DoDeviceEnrolmentCommandMethod = "DoDeviceEnrolment"
	GenerateVPCommandMethod        = "GenerateVP"
	GetVCredentialCommandMethod    = "GetVCredential"
	AcceptEnrolmentCommandMethod   = "AcceptEnrolment"
	VerifyCredentialCommandMethod  = "Verify"
	TestingCallMethod              = "TestingCall"
	GetTrustedIssuerListMethod     = "GetTrustedIssuerList"
	SignJWTContentCommandMethod    = "SignJWTContent"
	VerifyJWTContentCommandMethod  = "VerifyJWTContent"
	StoreCredentialCommandMethod   = "StoreCredential"
	DeriveProofCommandMethod       = "DeriveProof"
	// error messages.
	errEmptyNewDID       = "keys is mandatory"
	errEmptyUrl          = "url is mandatory"
	errEmptyDID          = "theirDid is mandatory"
	errEmptyIdProofs     = "idProofs is mandatory"
	erremptyCredId       = "credId is mandatory"
	errEmptyQueryByFrame = "querybyframe is mandatory"
	errEmptyContent      = "Content is mandatory"
	errEmptyJWT          = "JWT is mandatory"

	// log constants.
	didID = "did"

	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// provider contains dependencies for the vdr controller command operations
// and is typically created by using aries.Context().
type Provider interface {
	StorageProvider() storage.Provider
	VDRegistry() vdr.Registry
	Crypto() crypto.Crypto
	JSONLDDocumentLoader() ld.DocumentLoader
	MediaTypeProfiles() []string
}

// Command contains command operations provided by vdr controller.
type Command struct {
	vdrcommand        *vdrc.Command
	vcwalletcommand   *vcwalletc.Command
	walletuid         string
	walletpass        string
	currentDID        string //TODO UMU For retrieval of device DIDdoc, think about better implementation
	currentKeyPair    vcwalletc.CreateKeyPairResponse
	idProofValidators []IdProofValidator
	ctx               Provider
}

var doenrolmentMem = uint64(0)
var generateVPMem = uint64(0)
var verifyMem = uint64(0)

// New returns new poc client controller command instance.
func New(vdrcommand *vdrc.Command, vcwalletcommand *vcwalletc.Command) (*Command, error) {
	var idProofValidators []IdProofValidator
	idProofValidators = append(idProofValidators)

	//TODO UMU Add array (ordered) of validators and add validators for PoC
	idProofValidators = append(idProofValidators, &DefaultValidator{})

	src := rand.NewSource(time.Now().UnixNano())
	n := 12
	uid := randStringBytesMaskImprSrcUnsafe(n, src)
	pass := randStringBytesMaskImprSrcUnsafe(n, src)

	logutil.LogInfo(logger, "poc", "New", "uid: "+uid+" pass: "+pass)

	var l bytes.Buffer
	reader, err := getReader(&vcwalletc.CreateOrUpdateProfileRequest{
		UserID:             uid,
		LocalKMSPassphrase: pass,
	})
	if err != nil {
		return nil, err
	}
	cmdErr := vcwalletcommand.CreateProfile(&l, reader)
	if cmdErr != nil {
		return nil, cmdErr
	}
	return &Command{
		vdrcommand:        vdrcommand,
		vcwalletcommand:   vcwalletcommand,
		walletuid:         uid,
		walletpass:        pass,
		idProofValidators: idProofValidators,
	}, nil
}

// GetHandlers returns list of all commands supported by this controller command.
func (o *Command) GetHandlers() []command.Handler {
	return []command.Handler{
		cmdutil.NewCommandHandler(CommandName, NewDIDCommandMethod, o.NewDID),
		cmdutil.NewCommandHandler(CommandName, DoDeviceEnrolmentCommandMethod, o.DoDeviceEnrolment),
		cmdutil.NewCommandHandler(CommandName, GenerateVPCommandMethod, o.GenerateVP),
		cmdutil.NewCommandHandler(CommandName, AcceptEnrolmentCommandMethod, o.AcceptEnrolment),
		cmdutil.NewCommandHandler(CommandName, TestingCallMethod, o.TestingCall),
		cmdutil.NewCommandHandler(CommandName, VerifyCredentialCommandMethod, o.VerifyCredential),
		cmdutil.NewCommandHandler(CommandName, GetTrustedIssuerListMethod, o.GetTrustedIssuerList),
		cmdutil.NewCommandHandler(CommandName, StoreCredentialCommandMethod, o.StoreCredential),
		cmdutil.NewCommandHandler(CommandName, DeriveProofCommandMethod, o.DeriveProof),
	}
}

// testing call for memory usage
func (o *Command) TestingCall(rw io.Writer, req io.Reader) command.Error {

	//command.WriteNillableResponse(rw, &NewDIDResult{DIDDoc: parsedResponse.DID}, logger)

	//var err = o.vdrcommand.GetDID(&getResponse, reader)
	// var request TestingCallResult
	// err := json.NewDecoder(req).Decode(&request)
	// if err != nil {
	// 	logutil.LogInfo(logger, CommandName, AcceptEnrolmentCommandMethod, "failed to get DID: "+err.Error())
	// 	return command.NewValidationError(AcceptEnrolmentRequestErrorCode, fmt.Errorf("retrieve did doc error: %w", err))
	// }

	logutil.LogInfo(logger, CommandName, TestingCallMethod, "doenrolmentMem: "+strconv.FormatUint(doenrolmentMem, 10))
	logutil.LogInfo(logger, CommandName, TestingCallMethod, "generateVPMem: "+strconv.FormatUint(generateVPMem, 10))
	logutil.LogInfo(logger, CommandName, TestingCallMethod, "verifyMem: "+strconv.FormatUint(verifyMem, 10))
	testingCallResult := TestingCallResult{
		DoenrolmentMem: doenrolmentMem,
		GenerateVPMem:  generateVPMem,
		VerifyMem:      verifyMem,
	}
	logutil.LogInfo(logger, CommandName, TestingCallMethod, "example : "+strconv.FormatUint(testingCallResult.DoenrolmentMem, 10))
	command.WriteNillableResponse(rw, &TestingCallResult{DoenrolmentMem: doenrolmentMem, GenerateVPMem: generateVPMem, VerifyMem: verifyMem}, logger)
	logutil.LogInfo(logger, CommandName, TestingCallMethod, "success")
	return nil
}

// NewDID Generate and register DID for a set of new keys
func (o *Command) NewDID(rw io.Writer, req io.Reader) command.Error {
	var request NewDIDArgs

	err := json.NewDecoder(req).Decode(&request)
	if err != nil {
		logutil.LogInfo(logger, CommandName, NewDIDCommandMethod, err.Error())
		return command.NewValidationError(InvalidRequestErrorCode, fmt.Errorf("request decode : %w", err))
	}

	if request.Keys == nil || !checkAuthKeyPresent(request.Keys) {
		logutil.LogInfo(logger, CommandName, NewDIDCommandMethod, errEmptyNewDID)
		return command.NewValidationError(InvalidRequestErrorCode, fmt.Errorf(errEmptyNewDID))
	}

	doc := did.Doc{}
	doc.Context = []string{"https://w3id.org/did/v1"}

	//Open wallet
	var l bytes.Buffer
	reader, err := getReader(&vcwalletc.UnlockWalletRequest{
		UserID:             o.walletuid,
		LocalKMSPassphrase: o.walletpass,
	})
	if err != nil {
		return command.NewValidationError(NewDIDRequestErrorCode, fmt.Errorf("open wallet error: %w", err))
	}
	err = o.vcwalletcommand.Open(&l, reader)
	if err != nil {
		return command.NewValidationError(NewDIDRequestErrorCode, fmt.Errorf("open wallet error: %w", err))
	}
	token := getUnlockToken(l)
	if token == "" {
		return command.NewValidationError(NewDIDRequestErrorCode, fmt.Errorf("open wallet error decoding token"))
	}
	//Defer close wallet
	defer func() {
		var l2 bytes.Buffer
		reader, err = getReader(&vcwalletc.LockWalletRequest{
			UserID: o.walletuid,
		})
		err = o.vcwalletcommand.Close(&l2, reader)
		//TODO UMU See how to treat errors in this case
	}()

	//Compute id
	for idx, keyPurpose := range request.Keys {
		kt := parseKeyType(keyPurpose.KeyType)
		if kt == "" {
			logutil.LogInfo(logger, CommandName, NewDIDCommandMethod, "invalid key type")
			return command.NewValidationError(InvalidRequestErrorCode, fmt.Errorf("invalid key type"))
		}
		//parse number of keypurpose.keytype.Attrs for increment in 1
		if len(keyPurpose.KeyType.Attrs) > 0 {
			nAttrs := keyPurpose.KeyType.Attrs[0]
			nAug, err := strconv.Atoi(nAttrs)
			if err != nil {
				logutil.LogInfo(logger, CommandName, NewDIDCommandMethod, "parse number of key purpose key type attrs error")
				return command.NewValidationError(NewDIDRequestErrorCode, fmt.Errorf("parse number of key purpose key type attrs error: %w", err))
			}
			nAug = nAug + 1
			newAttrNumber := strconv.Itoa(nAug)
			keyPurpose.KeyType.Attrs[0] = newAttrNumber
		}
		reader, err = getReader(&vcwalletc.CreateKeyPairRequest{
			KeyType:    kt,
			WalletAuth: vcwalletc.WalletAuth{UserID: o.walletuid, Auth: token},
			Attrs:      keyPurpose.KeyType.Attrs,
		})
		var getResponse bytes.Buffer
		err = o.vcwalletcommand.CreateKeyPair(&getResponse, reader)
		if err != nil {
			return command.NewValidationError(NewDIDRequestErrorCode, fmt.Errorf("create key pair error: %w", err))
		}
		var parsedResponse vcwalletc.CreateKeyPairResponse
		err = json.NewDecoder(&getResponse).Decode(&parsedResponse)
		if err != nil {
			return command.NewValidationError(NewDIDRequestErrorCode, fmt.Errorf("create key pair error: %w", err))
		}

		keyID := parsedResponse.KeyID
		publicKeyb64 := parsedResponse.PublicKey
		if idx == 0 {
			doc.ID = "did:fabric:" + publicKeyb64
		}
		rawKey, err := base64.RawURLEncoding.DecodeString(publicKeyb64)
		if err != nil {
			logutil.LogInfo(logger, CommandName, NewDIDCommandMethod, fmt.Sprintf("parse b64 key error: request: %v: response: %v %v", string(kt), parsedResponse.KeyID, parsedResponse.PublicKey))
			return command.NewValidationError(NewDIDRequestErrorCode, fmt.Errorf("parse b64 key error: %w", err))
		}

		docKeyID := doc.ID + "#" + keyID

		verificationMethod := did.VerificationMethod{
			ID:         docKeyID,
			Type:       keyPurpose.KeyType.Type,
			Controller: doc.ID,
			Value:      rawKey,
		}

		doc.VerificationMethod = append(doc.VerificationMethod, verificationMethod)

		switch keyPurpose.Purpose {
		case "AssertionMethod":
			doc.AssertionMethod = append(doc.AssertionMethod, did.Verification{VerificationMethod: verificationMethod,
				Relationship: did.AssertionMethod})
		case "KeyAgreement":
			doc.KeyAgreement = append(doc.KeyAgreement, did.Verification{VerificationMethod: verificationMethod,
				Relationship: did.KeyAgreement})
		case "Authentication":
			doc.Authentication = append(doc.Authentication, did.Verification{VerificationMethod: verificationMethod,
				Relationship: did.Authentication})
			o.currentKeyPair = parsedResponse
		case "CapabilityDelegation":
			doc.CapabilityDelegation = append(doc.CapabilityDelegation, did.Verification{VerificationMethod: verificationMethod,
				Relationship: did.CapabilityDelegation})
		case "CapabilityInvocation":
			doc.CapabilityInvocation = append(doc.CapabilityInvocation, did.Verification{VerificationMethod: verificationMethod,
				Relationship: did.CapabilityInvocation})
		default: //If nothing we assume authentication
			doc.Authentication = append(doc.AssertionMethod, did.Verification{VerificationMethod: verificationMethod,
				Relationship: did.Authentication})
			o.currentKeyPair = parsedResponse
		}
	}
	now := time.Now()
	doc.Created = &now
	//Create DID
	var l1 bytes.Buffer
	marshalDocRequest, err := doc.JSONBytes()

	other := string(marshalDocRequest[:])
	fmt.Println(other)
	if err != nil {
		logutil.LogInfo(logger, CommandName, NewDIDCommandMethod, "failed to marshal DID Doc Request: "+err.Error())
		return command.NewValidationError(NewDIDRequestErrorCode, fmt.Errorf("marshalling did document request: %w", err))
	}
	opts := make(map[string]interface{})
	reader, err = getReader(&vdrc.CreateDIDRequest{
		Method: "fabric",
		DID:    marshalDocRequest,
		Opts:   opts,
	})
	err = o.vdrcommand.CreateDID(&l1, reader)
	if err != nil {
		return command.NewValidationError(NewDIDRequestErrorCode, fmt.Errorf("did creation error: %w", err))
	}
	var parsedResponse vdrc.Document
	err = json.NewDecoder(&l1).Decode(&parsedResponse)
	if err != nil {
		logutil.LogInfo(logger, CommandName, NewDIDCommandMethod, "failed to decode DID Document: "+err.Error())
		return command.NewValidationError(NewDIDRequestErrorCode, fmt.Errorf("did creation response error: %w", err))
	}
	o.currentDID = getDID(parsedResponse)
	if o.currentDID == "" {
		logutil.LogInfo(logger, CommandName, NewDIDCommandMethod, "get did error: (empty did)")
		return command.NewValidationError(NewDIDRequestErrorCode, fmt.Errorf("failed to parse id for future retrieval of document: %w", err))
	}
	//Save DID
	var l11 bytes.Buffer
	reader, err = getReader(&vdrc.DIDArgs{
		Document: parsedResponse,
		Name:     request.Name,
	})
	if err != nil {
		return command.NewValidationError(NewDIDRequestErrorCode, fmt.Errorf("save did error: %w", err))
	}
	err = o.vdrcommand.SaveDID(&l11, reader)
	if err != nil {
		return command.NewValidationError(NewDIDRequestErrorCode, fmt.Errorf("save did error: %w", err))
	}
	// finished
	command.WriteNillableResponse(rw, &NewDIDResult{DIDDoc: parsedResponse.DID}, logger)
	logutil.LogInfo(logger, CommandName, NewDIDCommandMethod, "success")
	//testing
	return nil
}

func (o *Command) SignJWTContent(rw io.Writer, req io.Reader) command.Error {
	var request SignJWTContentArgs

	err := json.NewDecoder(req).Decode(&request)
	if err != nil {
		logutil.LogInfo(logger, CommandName, SignJWTContentCommandMethod, err.Error())
		return command.NewValidationError(InvalidRequestErrorCode, fmt.Errorf("request decode : %w", err))
	}

	if request.Content == nil {
		logutil.LogInfo(logger, CommandName, SignJWTContentCommandMethod, errEmptyContent)
		return command.NewValidationError(InvalidRequestErrorCode, fmt.Errorf(errEmptyContent))
	}

	//Open wallet
	var l bytes.Buffer
	reader, err := getReader(&vcwalletc.UnlockWalletRequest{
		UserID:             o.walletuid,
		LocalKMSPassphrase: o.walletpass,
	})
	if err != nil {
		return command.NewValidationError(NewDIDRequestErrorCode, fmt.Errorf("open wallet error: %w", err))
	}
	err = o.vcwalletcommand.Open(&l, reader)
	if err != nil {
		return command.NewValidationError(NewDIDRequestErrorCode, fmt.Errorf("open wallet error: %w", err))
	}
	token := getUnlockToken(l)
	if token == "" {
		return command.NewValidationError(NewDIDRequestErrorCode, fmt.Errorf("open wallet error decoding token"))
	}
	//Defer close wallet
	defer func() {
		var l2 bytes.Buffer
		reader, err = getReader(&vcwalletc.LockWalletRequest{
			UserID: o.walletuid,
		})
		err = o.vcwalletcommand.Close(&l2, reader)
		//TODO UMU See how to treat errors in this case
	}()

	var content map[string]interface{}

	// Unmarshal the json.RawMessage into the map
	errUnmarshal := json.Unmarshal(request.Content, &content)
	if errUnmarshal != nil {
		fmt.Println("Error unmarshalling json.RawMessage:", err)
		return command.NewValidationError(SignJWTContentErrorCode, fmt.Errorf("error unmarshalling json.RawMessage: %w", err))
	}

	reqJWT := vcwalletc.SignJWTRequest{
		WalletAuth: vcwalletc.WalletAuth{UserID: o.walletuid, Auth: token},
		Headers:    nil,
		Claims:     content,
		KID:        o.currentDID + "#" + o.currentKeyPair.KeyID,
	}

	reqData, err := json.Marshal(reqJWT)
	if err != nil {
		logutil.LogInfo(logger, CommandName, "SignJWT", "failed to marshal request: "+err.Error())
	}
	requestFormatted := bytes.NewReader(reqData)
	// Capture the output
	var signBuf bytes.Buffer

	// Sign the JWT
	if err := o.vcwalletcommand.SignJWT(&signBuf, requestFormatted); err != nil {
		logutil.LogInfo(logger, CommandName, "SignJWT", "failed to sign JWT: "+err.Error())
	}

	var jwtResponse vcwalletc.SignJWTResponse

	err = json.Unmarshal(signBuf.Bytes(), &jwtResponse)
	if err != nil {
		logutil.LogInfo(logger, CommandName, "SignJWT", "failed to unmarshal JWT: "+err.Error())
	}

	signedJWT := jwtResponse.JWT
	fmt.Println("Signed JWT:", signedJWT)
	//Write the signedJWT as response
	command.WriteNillableResponse(rw, &SignJWTContentResult{SignedJWTContent: signedJWT}, logger)

	return nil
}

func (o *Command) VerifyJWTContent(rw io.Writer, req io.Reader) command.Error {

	var request VerifyJWTContentArgs

	err := json.NewDecoder(req).Decode(&request)
	if err != nil {
		logutil.LogInfo(logger, CommandName, VerifyJWTContentCommandMethod, err.Error())
		return command.NewValidationError(InvalidRequestErrorCode, fmt.Errorf("request decode : %w", err))
	}

	if request.JWT == "" {
		logutil.LogInfo(logger, CommandName, VerifyJWTContentCommandMethod, errEmptyContent)
		return command.NewValidationError(InvalidRequestErrorCode, fmt.Errorf(errEmptyJWT))
	}

	//Open wallet
	var l bytes.Buffer
	reader, err := getReader(&vcwalletc.UnlockWalletRequest{
		UserID:             o.walletuid,
		LocalKMSPassphrase: o.walletpass,
	})
	if err != nil {
		return command.NewValidationError(NewDIDRequestErrorCode, fmt.Errorf("open wallet error: %w", err))
	}
	err = o.vcwalletcommand.Open(&l, reader)
	if err != nil {
		return command.NewValidationError(NewDIDRequestErrorCode, fmt.Errorf("open wallet error: %w", err))
	}
	token := getUnlockToken(l)
	if token == "" {
		return command.NewValidationError(NewDIDRequestErrorCode, fmt.Errorf("open wallet error decoding token"))
	}
	//Defer close wallet
	defer func() {
		var l2 bytes.Buffer
		reader, err = getReader(&vcwalletc.LockWalletRequest{
			UserID: o.walletuid,
		})
		err = o.vcwalletcommand.Close(&l2, reader)
		//TODO UMU See how to treat errors in this case
	}()

	// Verify JWT
	verifyReq := &vcwalletc.VerifyJWTRequest{
		WalletAuth: vcwalletc.WalletAuth{UserID: o.walletuid, Auth: token},
		JWT:        request.JWT,
	}

	verifyReqBytes, _ := json.Marshal(verifyReq)
	verifyReqReader := bytes.NewReader(verifyReqBytes)
	var verifyBuf bytes.Buffer

	errVerify := o.vcwalletcommand.VerifyJWT(&verifyBuf, verifyReqReader)
	if errVerify != nil {
		logutil.LogInfo(logger, CommandName, VerifyJWTContentCommandMethod, "failed to verify JWT: "+err.Error())
	}
	fmt.Println("Verification result:", verifyBuf.String())
	//wrapp verifyBuf in VerifyJWTResponse

	var jwtVerifyResponse vcwalletc.VerifyJWTResponse

	errResp := json.Unmarshal(verifyBuf.Bytes(), &jwtVerifyResponse)
	if errResp != nil {
		logutil.LogInfo(logger, CommandName, "VerifyJWT", "failed to unmarshal JWT Verify Response: "+err.Error())
	}

	//write the verifyjwtresponse as response
	command.WriteNillableResponse(rw, jwtVerifyResponse, logger)

	return nil

}

func (o *Command) signJWT(token string) string {

	request := vcwalletc.SignJWTRequest{
		WalletAuth: vcwalletc.WalletAuth{UserID: o.walletuid, Auth: token},
		Headers:    nil,
		Claims: map[string]interface{}{
			"attrName":  "DID",
			"attrValue": o.currentDID,
		},
		KID: o.currentDID + "#" + o.currentKeyPair.KeyID,
	}

	reqData, err := json.Marshal(request)
	if err != nil {
		logutil.LogInfo(logger, CommandName, "SignJWT", "failed to marshal request: "+err.Error())
	}
	req := bytes.NewReader(reqData)
	// Capture the output
	var signBuf bytes.Buffer

	// Sign the JWT
	if err := o.vcwalletcommand.SignJWT(&signBuf, req); err != nil {
		logutil.LogInfo(logger, CommandName, "SignJWT", "failed to sign JWT: "+err.Error())
	}

	var jwtResponse vcwalletc.SignJWTResponse

	err = json.Unmarshal(signBuf.Bytes(), &jwtResponse)
	if err != nil {
		logutil.LogInfo(logger, CommandName, "SignJWT", "failed to unmarshal JWT: "+err.Error())
	}

	signedJWT := jwtResponse.JWT
	fmt.Println("Signed JWT:", signedJWT)
	return signedJWT
}

// verifyJWT
func (o *Command) verifyJWT(token string, signedJWT string) bool {

	// Verify JWT
	verifyReq := &vcwalletc.VerifyJWTRequest{
		WalletAuth: vcwalletc.WalletAuth{UserID: o.walletuid, Auth: token},
		JWT:        signedJWT,
	}

	verifyReqBytes, _ := json.Marshal(verifyReq)
	verifyReqReader := bytes.NewReader(verifyReqBytes)
	var verifyBuf bytes.Buffer

	err := o.vcwalletcommand.VerifyJWT(&verifyBuf, verifyReqReader)
	if err != nil {
		logutil.LogInfo(logger, CommandName, "SignJWT", "failed to verify JWT: "+err.Error())
	}
	fmt.Println("Verification result:", verifyBuf.String())
	//wrapp verifyBuf in VerifyJWTResponse

	var jwtVerifyResponse vcwalletc.VerifyJWTResponse

	errResp := json.Unmarshal(verifyBuf.Bytes(), &jwtVerifyResponse)
	if errResp != nil {
		logutil.LogInfo(logger, CommandName, "VerifyJWT", "failed to unmarshal JWT Verify Response: "+err.Error())
	}

	isVerified := jwtVerifyResponse.Verified
	return isVerified

}

// DoDeviceEnrolment Device completes an enrolment process against an issuer
func (o *Command) DoDeviceEnrolment(rw io.Writer, req io.Reader) command.Error {
	//Parse request
	var request DoDeviceEnrolmentArgs
	err := json.NewDecoder(req).Decode(&request)
	if err != nil {
		logutil.LogInfo(logger, CommandName, DoDeviceEnrolmentCommandMethod, err.Error())
		return command.NewValidationError(InvalidRequestErrorCode, fmt.Errorf("request decode : %w", err))
	}

	if request.Url == "" {
		logutil.LogInfo(logger, CommandName, DoDeviceEnrolmentCommandMethod, errEmptyUrl)
		return command.NewValidationError(InvalidRequestErrorCode, fmt.Errorf(errEmptyUrl))
	}

	//theirDID is not mandatory anymore
	// if request.TheirDID == "" {
	// 	logutil.LogInfo(logger, CommandName, DoDeviceEnrolmentCommandMethod, errEmptyDID)
	// 	return command.NewValidationError(InvalidRequestErrorCode, fmt.Errorf(errEmptyDID))
	// }

	if request.IdProofs == nil {
		logutil.LogInfo(logger, CommandName, DoDeviceEnrolmentCommandMethod, errEmptyIdProofs)
		return command.NewValidationError(InvalidRequestErrorCode, fmt.Errorf(errEmptyIdProofs))
	}

	identityProods := request.IdProofs

	//add current did to idProofs and sign with DID proofData with signJWT function

	//Open wallet
	var l bytes.Buffer
	reader, err := getReader(&vcwalletc.UnlockWalletRequest{
		UserID:             o.walletuid,
		LocalKMSPassphrase: o.walletpass,
	})
	if err != nil {
		return command.NewValidationError(DoDeviceEnrolmentRequestErrorCode, fmt.Errorf("open wallet error: %w", err))
	}
	err = o.vcwalletcommand.Open(&l, reader)
	if err != nil {
		return command.NewValidationError(DoDeviceEnrolmentRequestErrorCode, fmt.Errorf("open wallet error: %w", err))
	}
	token := getUnlockToken(l)
	if token == "" {
		logutil.LogInfo(logger, CommandName, DoDeviceEnrolmentCommandMethod, "could not get unlock token (empty token)")
		return command.NewValidationError(DoDeviceEnrolmentRequestErrorCode, fmt.Errorf("open wallet error decoding token"))
	}
	//Defer close wallet
	defer func() {
		var l2 bytes.Buffer
		reader, err = getReader(&vcwalletc.LockWalletRequest{
			UserID: o.walletuid,
		})
		err = o.vcwalletcommand.Close(&l2, reader)
		//TODO UMU See how to treat errors in this case
	}()

	//proofData := o.signJWT(token)
	//proofDataBytes := json.RawMessage(proofData)
	identityProods = append(identityProods, IdProof{AttrName: "DID", AttrValue: o.currentDID})

	// Do a post for AcceptEnrolmentResult to specified url
	acceptEnrolmentRequest := AcceptEnrolmentArgs{IdProofs: identityProods}
	jsonBody, err := json.Marshal(acceptEnrolmentRequest)

	if err != nil {
		logutil.LogInfo(logger, CommandName, DoDeviceEnrolmentCommandMethod, "could not generate request body")
		return command.NewValidationError(DoDeviceEnrolmentRequestErrorCode, fmt.Errorf("could not generate request body: %w", err))
	}

	//testing https insecure(for poc at the moment)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Post(request.Url+"/fluidos/idm/acceptEnrolment", "application/json", bytes.NewBuffer(jsonBody))

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	doenrolmentMem = m.Sys

	if err != nil {
		logutil.LogInfo(logger, CommandName, DoDeviceEnrolmentCommandMethod, "could not complete AcceptEnrolment POST request")
		return command.NewValidationError(DoDeviceEnrolmentRequestErrorCode, fmt.Errorf("could not complete AcceptEnrolment POST request: %w", err))
	}
	var res AcceptEnrolmentResult
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		logutil.LogInfo(logger, CommandName, DoDeviceEnrolmentCommandMethod, "could not parse AcceptEnrolment POST result")
		return command.NewValidationError(DoDeviceEnrolmentRequestErrorCode, fmt.Errorf("could not parse AcceptEnrolment POST result: %w", err))
	}
	if len(res.Credential) == 0 { //TODO UMU Better error message
		logutil.LogInfo(logger, CommandName, DoDeviceEnrolmentCommandMethod, "credential issuance was not completed")
		return command.NewValidationError(DoDeviceEnrolmentRequestErrorCode, fmt.Errorf("credential issuance was not completed: %s", res))
	}

	//Store cred in wallet
	serialCred, err := res.Credential.MarshalJSON()

	if err != nil {
		logutil.LogInfo(logger, CommandName, DoDeviceEnrolmentCommandMethod, "could not serialize cred")
		return command.NewValidationError(DoDeviceEnrolmentRequestErrorCode, fmt.Errorf("could not serialize cred: %w", err))
	}
	logutil.LogInfo(logger, CommandName, DoDeviceEnrolmentCommandMethod, "credential", string(serialCred))

	reader, err = getReader(&vcwalletc.AddContentRequest{
		WalletAuth:   vcwalletc.WalletAuth{UserID: o.walletuid, Auth: token},
		ContentType:  wallet.Credential,
		Content:      serialCred,
		CollectionID: "",
	})
	var getResponse bytes.Buffer
	err = o.vcwalletcommand.Add(&getResponse, reader)
	if err != nil {
		logutil.LogInfo(logger, CommandName, DoDeviceEnrolmentCommandMethod, "store credential error")
		return command.NewValidationError(DoDeviceEnrolmentRequestErrorCode, fmt.Errorf("store credential error: %w", err))
	}
	//Return cred
	stoId, err := getContentID(serialCred)
	if err != nil {
		logutil.LogInfo(logger, CommandName, DoDeviceEnrolmentCommandMethod, "storage id error")
		return command.NewValidationError(DoDeviceEnrolmentRequestErrorCode, fmt.Errorf("storage id error: %w", err))
	}
	command.WriteNillableResponse(rw, &DoDeviceEnrolmentResult{Credential: res.Credential, CredStorageId: stoId}, logger)
	logutil.LogInfo(logger, CommandName, DoDeviceEnrolmentCommandMethod, "success")
	return nil
}

func (o *Command) GetVCredential(rw io.Writer, req io.Reader) command.Error {
	var request GetVCredentialArgs
	err := json.NewDecoder(req).Decode(&request)
	if err != nil {
		logutil.LogInfo(logger, CommandName, GetVCredentialCommandMethod, err.Error())
		return command.NewValidationError(InvalidRequestErrorCode, fmt.Errorf("request decode : %w", err))
	}
	if request.CredId == "" {
		logutil.LogInfo(logger, CommandName, GetVCredentialCommandMethod, erremptyCredId)
		return command.NewValidationError(InvalidRequestErrorCode, fmt.Errorf(erremptyCredId))
	}
	//Open wallet
	var l bytes.Buffer
	reader, err := getReader(&vcwalletc.UnlockWalletRequest{
		UserID:             o.walletuid,
		LocalKMSPassphrase: o.walletpass,
	})
	if err != nil {
		return command.NewValidationError(GenerateVPRequestErrorCode, fmt.Errorf("open wallet error: %w", err))
	}
	err = o.vcwalletcommand.Open(&l, reader)
	if err != nil {
		return command.NewValidationError(GenerateVPRequestErrorCode, fmt.Errorf("open wallet error: %w", err))
	}
	token := getUnlockToken(l)
	if token == "" {
		logutil.LogInfo(logger, CommandName, GenerateVPCommandMethod, "failed to get unlock token (empty token)")
		return command.NewValidationError(GenerateVPRequestErrorCode, fmt.Errorf("open wallet error decoding token"))
	}
	//Defer close wallet
	defer func() {
		var l2 bytes.Buffer
		reader, err = getReader(&vcwalletc.LockWalletRequest{
			UserID: o.walletuid,
		})
		err = o.vcwalletcommand.Close(&l2, reader)
	}()
	//Get stored credential from Id
	//var credID = request.CredId
	reader, err = getReader(&vcwalletc.GetContentRequest{
		ContentID:   request.CredId,
		ContentType: wallet.Credential,
		WalletAuth:  vcwalletc.WalletAuth{UserID: o.walletuid, Auth: token},
	})

	var getResponse bytes.Buffer
	err = o.vcwalletcommand.Get(&getResponse, reader)
	if err != nil {
		return command.NewValidationError(GenerateVPRequestErrorCode, fmt.Errorf("retrieve credential error: %w", err))
	}
	var parsedResponse vcwalletc.GetContentResponse
	err = json.NewDecoder(&getResponse).Decode(&parsedResponse)
	if err != nil {
		return command.NewValidationError(GenerateVPRequestErrorCode, fmt.Errorf("retrieve credential error: %w", err))
	}

	if err != nil {
		return command.NewValidationError(GenerateVPRequestErrorCode, fmt.Errorf("failed to decode stored credential: %w", err))
	}

	command.WriteNillableResponse(rw, &GetVCredentialResult{parsedResponse.Content}, logger)
	return nil
}

// GenerateVP Device generates VPresentation (or VCredential for now) for an authorization process
func (o *Command) GenerateVP(rw io.Writer, req io.Reader) command.Error {
	//TODO UMU For now we use ContentId, but we should do it through query or similar and might even be simpler
	//XXX UMU  (maybe Query would be useful for real implementation of credential retrieval). If VPresentation, it might be signed with DIDKey or something like that
	//Parse parameters from request
	var request GenerateVPArgs
	err := json.NewDecoder(req).Decode(&request)
	if err != nil {
		logutil.LogInfo(logger, CommandName, GenerateVPCommandMethod, err.Error())
		return command.NewValidationError(InvalidRequestErrorCode, fmt.Errorf("request decode : %w", err))
	}
	// if request.CredId == "" {
	// 	logutil.LogInfo(logger, CommandName, GenerateVPCommandMethod, errEmptyUrl)
	// 	return command.NewValidationError(InvalidRequestErrorCode, fmt.Errorf(erremptyCredId))
	// }
	// if request.QueryByFrame == nil {
	// 	logutil.LogInfo(logger, CommandName, GenerateVPCommandMethod, errEmptyUrl)
	// 	return command.NewValidationError(InvalidRequestErrorCode, fmt.Errorf(errEmptyQueryByFrame))
	// }

	//Open wallet
	var l bytes.Buffer
	reader, err := getReader(&vcwalletc.UnlockWalletRequest{
		UserID:             o.walletuid,
		LocalKMSPassphrase: o.walletpass,
	})
	if err != nil {
		return command.NewValidationError(GenerateVPRequestErrorCode, fmt.Errorf("open wallet error: %w", err))
	}
	err = o.vcwalletcommand.Open(&l, reader)
	if err != nil {
		return command.NewValidationError(GenerateVPRequestErrorCode, fmt.Errorf("open wallet error: %w", err))
	}
	token := getUnlockToken(l)
	if token == "" {
		logutil.LogInfo(logger, CommandName, GenerateVPCommandMethod, "failed to get unlock token (empty token)")
		return command.NewValidationError(GenerateVPRequestErrorCode, fmt.Errorf("open wallet error decoding token"))
	}
	//Defer close wallet
	defer func() {
		var l2 bytes.Buffer
		reader, err = getReader(&vcwalletc.LockWalletRequest{
			UserID: o.walletuid,
		})
		err = o.vcwalletcommand.Close(&l2, reader)
	}()
	//Get stored credential from Id
	//var credID = request.CredId
	reader, err = getReader(&vcwalletc.GetContentRequest{
		ContentID:   request.CredId,
		ContentType: wallet.Credential,
		WalletAuth:  vcwalletc.WalletAuth{UserID: o.walletuid, Auth: token},
	})

	var getResponse bytes.Buffer
	err = o.vcwalletcommand.Get(&getResponse, reader)
	if err != nil {
		return command.NewValidationError(GenerateVPRequestErrorCode, fmt.Errorf("retrieve credential error: %w", err))
	}
	var parsedResponse vcwalletc.GetContentResponse
	err = json.NewDecoder(&getResponse).Decode(&parsedResponse)
	if err != nil {
		return command.NewValidationError(GenerateVPRequestErrorCode, fmt.Errorf("retrieve credential error: %w", err))
	}

	if err != nil {
		return command.NewValidationError(GenerateVPRequestErrorCode, fmt.Errorf("failed to decode stored credential: %w", err))
	}

	//treatment for selective disclosure frame (query by frame)
	var rawMessages []json.RawMessage
	frameBytes, err := json.Marshal(request.QueryByFrame)
	if err != nil {
		fmt.Println("Error marshaling Frame:", err)
	}
	rawMessages = append(rawMessages, json.RawMessage(frameBytes))

	reader, err = getReader(&vcwalletc.ContentQueryRequest{
		WalletAuth: vcwalletc.WalletAuth{UserID: o.walletuid, Auth: token},
		Query: []*wallet.QueryParams{
			{
				Type:  "QueryByFrame",
				Query: rawMessages,
			},
		},
	})

	var queryResponse bytes.Buffer
	queryErr := o.vcwalletcommand.Query(&queryResponse, reader)
	if queryErr != nil {
		return command.NewValidationError(GenerateVPRequestErrorCode, fmt.Errorf("query response not working: %w", queryErr))
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	generateVPMem = m.Sys

	var queryParsedResponse vcwalletc.CustomContentQueryResponse

	err = json.Unmarshal(queryResponse.Bytes(), &queryParsedResponse)
	if err != nil {
		return command.NewValidationError(GenerateVPRequestErrorCode, fmt.Errorf("unmarshal not working: %w", err))
	}
	logutil.LogInfo(logger, CommandName, GenerateVPCommandMethod, "Verifiable Presentation result response without unmarshall: "+queryResponse.String())

	command.WriteNillableResponse(rw, &GenerateVPResultCustom{queryParsedResponse.Results}, logger)

	logutil.LogInfo(logger, CommandName, GenerateVPCommandMethod, "success")

	return nil
}

func (o *Command) VerifyCredential(rw io.Writer, req io.Reader) command.Error {
	logutil.LogDebug(logger, CommandName, "validateCredential", "start")
	//TODO UMU: create method for command and rest
	var request VerifyCredentialArgs
	err := json.NewDecoder(req).Decode(&request)
	if err != nil {
		return command.NewValidationError(InvalidRequestErrorCode, fmt.Errorf("request decode : %w", err))
	}

	//Open wallet
	var l bytes.Buffer
	reader, err := getReader(&vcwalletc.UnlockWalletRequest{
		UserID:             o.walletuid,
		LocalKMSPassphrase: o.walletpass,
	})
	if err != nil {
		return command.NewValidationError(VerifyCredentialRequestErrorCode, fmt.Errorf("open wallet error: %w", err))
	}
	err = o.vcwalletcommand.Open(&l, reader)
	if err != nil {
		logutil.LogError(logger, CommandName, VerifyCredentialCommandMethod, "failed to open wallet: "+err.Error())
		return command.NewValidationError(VerifyCredentialRequestErrorCode, fmt.Errorf("open wallet error: %w", err))
	}
	//Defer close wallet
	defer func() {
		var l2 bytes.Buffer
		reader, err = getReader(&vcwalletc.LockWalletRequest{
			UserID: o.walletuid,
		})
		err = o.vcwalletcommand.Close(&l2, reader)
	}()

	token := getUnlockToken(l)
	if token == "" {
		logutil.LogInfo(logger, CommandName, VerifyCredentialCommandMethod, "failed to get unlock token (empty token)")
		return command.NewValidationError(VerifyCredentialRequestErrorCode, fmt.Errorf("open wallet error decoding token"))
	}

	if err != nil {
		logutil.LogError(logger, CommandName, VerifyCredentialCommandMethod, "failed to marshal credential: "+err.Error())
		return command.NewValidationError(VerifyCredentialRequestErrorCode, fmt.Errorf("failed to marshal credential: %w", err))
	}

	var response vcwalletc.VerifyResponse

	useCred := true
	if request.Credential == nil {
		useCred = false
		if request.Presentation == nil {
			logutil.LogError(logger, CommandName, VerifyCredentialCommandMethod, "must include credential or presentation in request")
			return command.NewValidationError(VerifyCredentialRequestErrorCode, fmt.Errorf("must include credential or presentation in request"))
		}
	}

	vcwalletcommandrequest := vcwalletc.VerifyRequest{
		WalletAuth: vcwalletc.WalletAuth{UserID: o.walletuid, Auth: token},
	}
	if useCred {
		vcwalletcommandrequest.RawCredential = request.Credential
		logutil.LogDebug(logger, CommandName, VerifyCredentialCommandMethod, "what am i verifying? "+string(vcwalletcommandrequest.RawCredential))
	} else {
		vcwalletcommandrequest.Presentation = request.Presentation
		logutil.LogDebug(logger, CommandName, VerifyCredentialCommandMethod, "what am i verifying? "+string(vcwalletcommandrequest.Presentation))

	}
	reader, err = getReader(&vcwalletcommandrequest)

	//golang find and replace char in string
	if err != nil {
		return command.NewValidationError(VerifyCredentialRequestErrorCode, fmt.Errorf("failed to get Verify Request reader: %w", err))
	}
	var l2 bytes.Buffer
	err = o.vcwalletcommand.Verify(&l2, reader)
	if err != nil {
		logutil.LogDebug(logger, CommandName, VerifyCredentialCommandMethod, "error in Verify")
		command.WriteNillableResponse(rw, &VerifyCredentialResult{Result: false, Error: err.Error()}, logger)
		return nil
	}
	err = json.NewDecoder(&l2).Decode(&response)
	if err != nil {
		logutil.LogDebug(logger, CommandName, VerifyCredentialCommandMethod, "error in Decode verify Response")
		return command.NewValidationError(VerifyCredentialRequestErrorCode, fmt.Errorf("failed to decode verify response: %w", err))
	}
	var result string
	if !response.Verified {
		result = "not verified"
		//return command.NewValidationError(VerifyCredentialRequestErrorCode, fmt.Errorf("failed to verify credential: %s", response.Error))
		logutil.LogDebug(logger, CommandName, VerifyCredentialCommandMethod, "credential verified response: "+result)
		command.WriteNillableResponse(rw, &VerifyCredentialResult{Result: false, Error: "not valid: " + response.Error}, logger)
		return nil
	}
	result = "verified"
	logutil.LogDebug(logger, CommandName, VerifyCredentialCommandMethod, "credential verified response: "+result)
	command.WriteNillableResponse(rw, &VerifyCredentialResult{Result: response.Verified}, logger)
	return nil
}

// AcceptEnrolment Issuer exposes this method to acacept enrolment processes, ending in credential issuance
func (o *Command) AcceptEnrolment(rw io.Writer, req io.Reader) command.Error {
	//Parse request arguments
	var request AcceptEnrolmentArgs
	err := json.NewDecoder(req).Decode(&request)
	if err != nil {
		logutil.LogInfo(logger, CommandName, AcceptEnrolmentCommandMethod, err.Error())
		return command.NewValidationError(InvalidRequestErrorCode, fmt.Errorf("request decode : %w", err))
	}

	if request.IdProofs == nil {
		logutil.LogInfo(logger, CommandName, AcceptEnrolmentCommandMethod, errEmptyIdProofs)
		return command.NewValidationError(InvalidRequestErrorCode, fmt.Errorf(errEmptyIdProofs))
	}
	//Open wallet
	var l bytes.Buffer
	reader, err := getReader(&vcwalletc.UnlockWalletRequest{
		UserID:             o.walletuid,
		LocalKMSPassphrase: o.walletpass,
	})
	if err != nil {
		return command.NewValidationError(GenerateVPRequestErrorCode, fmt.Errorf("open wallet error: %w", err))
	}
	err = o.vcwalletcommand.Open(&l, reader)
	if err != nil {
		return command.NewValidationError(GenerateVPRequestErrorCode, fmt.Errorf("open wallet error: %w", err))
	}
	token := getUnlockToken(l)
	if token == "" {
		logutil.LogInfo(logger, CommandName, AcceptEnrolmentCommandMethod, "could not get unlock token (empty token)")
		return command.NewValidationError(GenerateVPRequestErrorCode, fmt.Errorf("open wallet error decoding token"))
	}
	//Defer close wallet
	defer func() {
		var l2 bytes.Buffer
		reader, err = getReader(&vcwalletc.LockWalletRequest{
			UserID: o.walletuid,
		})
		err = o.vcwalletcommand.Close(&l2, reader)
	}()
	//Initialize credential for issuance
	baseCredString := "{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.w3.org/2018/credentials/examples/v1\",\"https://ssiproject.inf.um.es/security/psms/v1\",\"https://ssiproject.inf.um.es/poc/context/v1\"],\"type\":[\"VerifiableCredential\",\"FluidosCredential\"]}"
	var baseCred map[string]interface{}
	err = json.Unmarshal([]byte(baseCredString), &baseCred)
	if err != nil {
		logutil.LogInfo(logger, CommandName, AcceptEnrolmentCommandMethod, "failed to decode base cred")
		return command.NewValidationError(AcceptEnrolmentRequestErrorCode, fmt.Errorf("failed to decode base cred"))
	}
	//Validate IdProofs and generate credentialSubject from them
	credSubject, validIdProofs, err := o.generateCredentialSubject(request.IdProofs)
	if !validIdProofs {
		logutil.LogInfo(logger, CommandName, AcceptEnrolmentCommandMethod, "failed to validate identity proofs")
		//TODO UMU Write response indicating failed IdProof to client?
		return command.NewValidationError(AcceptEnrolmentRequestErrorCode, fmt.Errorf("failed to validate identity proofs: %w", err))
	}
	if err != nil {
		logutil.LogInfo(logger, CommandName, AcceptEnrolmentCommandMethod, "failed to parse identity proofs into credential subject")
		return command.NewValidationError(AcceptEnrolmentRequestErrorCode, fmt.Errorf("failed to parse identity proofs into credential subject %w", err))
	}
	//baseCred["credentialSubject"] = credSubject
	baseCred["credentialSubject"] = make(map[string]string, len(credSubject))

	for k, v := range credSubject {
		baseCred["credentialSubject"].(map[string]string)[k] = v.(string)
	}

	//Get DID/DIDDoc for specifying key, issuer...
	reader, err = getReader(&vdrc.IDArg{
		ID: o.currentDID,
	})
	var getResponse bytes.Buffer
	err = o.vdrcommand.GetDID(&getResponse, reader)
	if err != nil {
		logutil.LogInfo(logger, CommandName, AcceptEnrolmentCommandMethod, "failed to get DID: "+err.Error())
		return command.NewValidationError(AcceptEnrolmentRequestErrorCode, fmt.Errorf("retrieve did doc error: %w", err))
	}
	var parsedDoc vdrc.Document
	err = json.NewDecoder(&getResponse).Decode(&parsedDoc)
	if err != nil {
		logutil.LogInfo(logger, CommandName, AcceptEnrolmentCommandMethod, "failed to decode DID Document: "+err.Error())
		return command.NewValidationError(AcceptEnrolmentRequestErrorCode, fmt.Errorf("retrieve did doc error: %w", err))
	}
	didDoc, err := did.ParseDocument(parsedDoc.DID)
	if err != nil {
		logutil.LogInfo(logger, CommandName, AcceptEnrolmentCommandMethod, "failed to parse DID Document: "+err.Error())
		return command.NewValidationError(AcceptEnrolmentRequestErrorCode, fmt.Errorf("retrieve did doc error: %w", err))
	}
	//Generate credential metadata: issuanceDate, expirationDate, id, issuer,
	now := time.Now()
	baseCred["issuanceDate"] = now //TODO UMU Linkability issue
	baseCred["expirationDate"] = now.Add(100000000000000)
	rand.Seed(time.Now().UnixNano())
	baseCred["id"] = didDoc.ID + strconv.Itoa(rand.Intn(10000000)) //XXX UMU See whether it is really necessary and in that case which value
	baseCred["issuer"] = didDoc.ID
	//Generate request and call issue
	reqCred, err := json.Marshal(baseCred)
	if err != nil {
		logutil.LogInfo(logger, CommandName, AcceptEnrolmentCommandMethod, "failed to marshall credential request")
		return command.NewValidationError(AcceptEnrolmentRequestErrorCode, fmt.Errorf("failed to marshall credential request %w", err))
	}
	proofRepr := verifiable.SignatureProofValue
	logutil.LogInfo(logger, CommandName, AcceptEnrolmentCommandMethod, string(reqCred))
	issueRequest := &vcwalletc.IssueRequest{
		Credential: reqCred,
		WalletAuth: vcwalletc.WalletAuth{UserID: o.walletuid, Auth: token},
		ProofOptions: &wallet.ProofOptions{Controller: didDoc.ID, //TODO UMU It should be interesting to also specify verificationMethod, for now it just takes the assertionMethod associated to controller
			ProofRepresentation: &proofRepr,
			ProofType:           wallet.PsmsBlsSignature2022,
		},
	}
	reader, err = getReader(issueRequest)
	if err != nil {
		logutil.LogInfo(logger, CommandName, AcceptEnrolmentCommandMethod, "failed to generate request for issue command")
		return command.NewValidationError(AcceptEnrolmentRequestErrorCode, fmt.Errorf("failed to generate request for issue command: %w", err))
	}
	var issueResponse bytes.Buffer
	err = o.vcwalletcommand.Issue(&issueResponse, reader)

	if err != nil {
		return command.NewValidationError(AcceptEnrolmentRequestErrorCode, fmt.Errorf("issuance error: %w", err))
	}
	logutil.LogInfo(logger, CommandName, AcceptEnrolmentCommandMethod, string(issueResponse.Bytes()))
	var parsedResponse AcceptEnrolmentResult
	err = json.NewDecoder(&issueResponse).Decode(&parsedResponse)
	if err != nil {
		return command.NewValidationError(AcceptEnrolmentRequestErrorCode, fmt.Errorf("issuance error: %w", err))
	}
	//Return result
	command.WriteNillableResponse(rw, &parsedResponse, logger)
	logutil.LogInfo(logger, CommandName, AcceptEnrolmentCommandMethod, "success")
	return nil
}

// GetTrustedIssuerList returns the list of trusted issuers, mocked for now
func (o *Command) GetTrustedIssuerList(rw io.Writer, req io.Reader) command.Error {
	//TODO UMU: Implement
	trustedIssuer := TrustedIssuer{
		DID:       "did:fabric:zxdkpwDnu7ixBidF_I8sgMI6Q4St0t90HY-_JmlHZFI",
		IssuerUrl: "https://issuer:9082",
	}
	var trustedIssuerList []TrustedIssuer
	trustedIssuerList = append(trustedIssuerList, trustedIssuer)

	var trustedIssuerListResponse = GetTrustedIssuerListResult{
		TrustedIssuers: trustedIssuerList,
	}

	command.WriteNillableResponse(rw, &trustedIssuerListResponse, logger)
	return nil
}

// / NANCY quick test additions
func (o *Command) StoreCredential(rw io.Writer, req io.Reader) command.Error {
	logutil.LogDebug(logger, CommandName, "validateCredential", "start")
	var request StoreCredentialArgs
	err := json.NewDecoder(req).Decode(&request)
	if err != nil {
		return command.NewValidationError(InvalidRequestErrorCode, fmt.Errorf("request decode : %w", err))
	}

	//Open wallet
	var l bytes.Buffer
	reader, err := getReader(&vcwalletc.UnlockWalletRequest{
		UserID:             o.walletuid,
		LocalKMSPassphrase: o.walletpass,
	})
	if err != nil {
		return command.NewValidationError(StoreCredentialRequestErrorCode, fmt.Errorf("open wallet error: %w", err))
	}
	err = o.vcwalletcommand.Open(&l, reader)
	if err != nil {
		logutil.LogError(logger, CommandName, StoreCredentialCommandMethod, "failed to open wallet: "+err.Error())
		return command.NewValidationError(StoreCredentialRequestErrorCode, fmt.Errorf("open wallet error: %w", err))
	}
	//Defer close wallet
	defer func() {
		var l2 bytes.Buffer
		reader, err = getReader(&vcwalletc.LockWalletRequest{
			UserID: o.walletuid,
		})
		err = o.vcwalletcommand.Close(&l2, reader)
	}()

	token := getUnlockToken(l)
	if token == "" {
		logutil.LogInfo(logger, CommandName, StoreCredentialCommandMethod, "failed to get unlock token (empty token)")
		return command.NewValidationError(StoreCredentialRequestErrorCode, fmt.Errorf("open wallet error decoding token"))
	}

	if err != nil {
		logutil.LogError(logger, CommandName, StoreCredentialCommandMethod, "failed to marshal credential: "+err.Error())
		return command.NewExecuteError(StoreCredentialRequestErrorCode, fmt.Errorf("failed to marshal credential: %w", err))
	}
	reader, err = getReader(&vcwalletc.AddContentRequest{
		WalletAuth:  vcwalletc.WalletAuth{UserID: o.walletuid, Auth: token},
		ContentType: wallet.Credential,
		Content:     request.Credential,
	})
	//golang find and replace char in string
	if err != nil {
		return command.NewExecuteError(StoreCredentialRequestErrorCode, fmt.Errorf("failed to get Store Request reader: %w", err))
	}
	var l2 bytes.Buffer
	err = o.vcwalletcommand.Add(&l2, reader)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	verifyMem = m.Sys
	if err != nil {
		return command.NewExecuteError(StoreCredentialRequestErrorCode, fmt.Errorf("failed to store credential: %w", err))
	}
	var result string
	logutil.LogDebug(logger, CommandName, StoreCredentialCommandMethod, "credential store response:"+result)
	command.WriteNillableResponse(rw, &StoreCredentialResult{Result: true}, logger)
	return nil
}

func (o *Command) DeriveProof(rw io.Writer, req io.Reader) command.Error {
	//Parse parameters from request
	var request DeriveProofArgs
	err := json.NewDecoder(req).Decode(&request)
	if err != nil {
		logutil.LogInfo(logger, CommandName, DeriveProofCommandMethod, err.Error())
		return command.NewValidationError(InvalidRequestErrorCode, fmt.Errorf("request decode : %w", err))
	}
	// if request.CredId == "" {
	// 	logutil.LogInfo(logger, CommandName, GenerateVPCommandMethod, errEmptyUrl)
	// 	return command.NewValidationError(InvalidRequestErrorCode, fmt.Errorf(erremptyCredId))
	// }
	// if request.QueryByFrame == nil {
	// 	logutil.LogInfo(logger, CommandName, GenerateVPCommandMethod, errEmptyUrl)
	// 	return command.NewValidationError(InvalidRequestErrorCode, fmt.Errorf(errEmptyQueryByFrame))
	// }

	//Open wallet
	var l bytes.Buffer
	reader, err := getReader(&vcwalletc.UnlockWalletRequest{
		UserID:             o.walletuid,
		LocalKMSPassphrase: o.walletpass,
	})
	if err != nil {
		return command.NewValidationError(DeriveProofRequestErrorCode, fmt.Errorf("open wallet error: %w", err))
	}
	err = o.vcwalletcommand.Open(&l, reader)
	if err != nil {
		return command.NewValidationError(DeriveProofRequestErrorCode, fmt.Errorf("open wallet error: %w", err))
	}
	token := getUnlockToken(l)
	if token == "" {
		logutil.LogInfo(logger, CommandName, GenerateVPCommandMethod, "failed to get unlock token (empty token)")
		return command.NewValidationError(DeriveProofRequestErrorCode, fmt.Errorf("open wallet error decoding token"))
	}
	//Defer close wallet
	defer func() {
		var l2 bytes.Buffer
		reader, err = getReader(&vcwalletc.LockWalletRequest{
			UserID: o.walletuid,
		})
		err = o.vcwalletcommand.Close(&l2, reader)
	}()

	reader, err = getReader(&vcwalletc.DeriveRequest{
		WalletAuth:         vcwalletc.WalletAuth{UserID: o.walletuid, Auth: token},
		StoredCredentialID: request.CredId,
		DeriveOptions: &wallet.DeriveOptions{
			Frame: request.QueryFrame,
			Nonce: request.Nonce,
		},
	})

	a, _ := json.Marshal(request.QueryFrame)
	logutil.LogDebug(logger, CommandName, DeriveProofCommandMethod, "Derive for id "+request.CredId+" Nonce "+request.Nonce+" Frame "+string(a))

	var queryResponse bytes.Buffer
	queryErr := o.vcwalletcommand.Derive(&queryResponse, reader)
	if queryErr != nil {
		return command.NewValidationError(DeriveProofRequestErrorCode, fmt.Errorf("query response not working: %w", queryErr))
	}

	var queryParsedResponse DeriveProofWalletResult

	err = json.NewDecoder(&queryResponse).Decode(&queryParsedResponse)
	if err != nil {
		return command.NewValidationError(AcceptEnrolmentRequestErrorCode, fmt.Errorf("issuance error: %w", err))
	}
	//XXXX Return result
	//command.WriteNillableResponse(rw, &queryParsedResponse, logger)

	logutil.LogDebug(logger, CommandName, DeriveProofCommandMethod, "Returning cred "+string(queryParsedResponse.Credential))

	command.WriteNillableResponse(rw, &DeriveProofResult{Result: queryParsedResponse.Credential}, logger)

	logutil.LogInfo(logger, CommandName, DeriveProofCommandMethod, "success")

	return nil
}
