# Moded Aries-Framework-Go for p-ABC usage in 5g authentication experiments
The deployment will be native in each device, the following steps should be followed:

## Prerequisites
Installing python (e.g. 3.X), gcc, Cmake and make for compiling the C libraries for the p-ABC code.
Install Go1.19.X for the agent.

## Deployment process
First get the code:
```
git clone https://github.com/JesusGarciaRodriguez/aries-framework-5g.git
git submodule update --init --recursive
```

Compile the C library for the p-ABC crypto with the following command:
```
sh setupscript.sh -b -d <directory-aries-framework>/pkg/crypto/primitive/psms12381g1pub/clib
```

Configure keys and certificates:
1. Copy the keys folder to  <directory-aries-framework>/test/bdd/fixtures/
2. Install certificates 
```
sudo mkdir /usr/local/share/ca-certificates/ssi
sudo cp ../test/bdd/fixtures/keys/tls/* /usr/local/share/ca-certificates/ssi
sudo cp /usr/local/share/ca-certificates/ssi/ec-cacert.pem /usr/local/share/ca-certificates/ssi/ec-cacert.crt
sudo cp /usr/local/share/ca-certificates/ssi/ec-pubCert.pem /usr/local/share/ca-certificates/ssi/ec-pubCert.crt
update-ca-certificates
```

Building the agent (this will be the only one necessary to rerun when changes to the agent are done):
```
cd <directory-aries-framework>/cmd/aries-agent-rest
go mod tidy &&  go build -o ../../build/bin/aries-agent-rest main.go
```

Running:
```
cd <directory-aries-framework>/scripts
BIN_PATH=../build/bin bash run-5g-demo-device.sh
```

## Benchmark
A set of benchmarks for the crypto operations within the wallet (sign,derive,verify credentials) can be run:
```
cd <directory-aries-framework>/pkg/bench
BIN_PATH=../build/bin bash run-5g-demo-device.sh
go test -bench=. -benchtime=50x -benchmem 
go test -bench=BenchmarkKeygen -benchtime=100x (to specify a benchmark method)
```

## Testing
To do a quick test of the functionality, the following REST petitions can be sent:
StoreCredential:
```
curl --location 'http://localhost:8082/nancy/idm/storeCredential' \
--data-raw '{
   "credential": {
   "@context":[
      "https://www.w3.org/2018/credentials/v1",
      "https://ssiproject.inf.um.es/security/psms/v1",
      "https://nancy-identity/context/exampleContext/v1",
      "https://w3id.org/security/bbs/v1"
   ],
   "credentialSubject":{
      "acces5g":[
         "99940",
         "99941"
      ],
      "accesMec":true,
      "accessLocEng":false,
      "accessV2X":true,
      "id":"did:example:ebfeb1f712ebc6f1c276e12ec21",
      "nancyId":"zbASDASDjfasfasfasfiasfjasfjioasfjas0asjfebfeb1f712ebc6f1c276e12ec21"
   },
   "expirationDate":"2024-01-01T19:23:24Z",
   "id":"http://example/credentials/18723",
   "issuanceDate":"2023-07-08T19:23:24Z",
   "issuer":"did:erat:example",
   "proof":{
      "created":"2025-02-24T14:11:53.741729213+01:00",
      "proofPurpose":"assertionMethod",
      "proofValue":"BAN2su2p1rTxxFkCBTRQ4EllWb5qC7hOTE7NWn2KFY0iJnC-cyhzsRX0gkOKR880uxNvQh06_TDeUs4r9Q1edI9ZoA-IybQDaPZJS0I-jA1e0tabDkn8KeZWH2zchgxY6xG_okLCVLyr5_5ryvg4-bv08MXs2GBzOfvzkbeJ-JKjk_OZw0sIiQCoJYSLp1jeiBKzL0GFmNPcMbJxn2kCoAUpAS6AXkxBu0JRmGkb5Gc_6aTy67AKfoSZQoLSQqQmMgQM-T6oiZtvvyf-pwq2kpzbTdbvZv1Kh04IW1Nj4-Lc1wWabC3oOce3bFOXij5HR1IQVGhtzi3KGTNZ3ewJcxTrBBPP0drzycuWrPHu6gc5qGkNTjXosADl8zGY7VEGl_QECOXgAhJoT_2GJWLVEkb6g2c41U97zYeyQAlpFFZ3DqsW3HXFjuPzZ1V5zW_BuW8JiMEI1_C77wDmCLRWzNnH6YhMwE1eR6FODKGJlfPqMzQs3gxUNf2JajcJuU0I-hIwdahvdLWPY6He5kxa9Yydft3CDpaq3NVjqfgrirWs9YLYQje-Zamgf_BvElIk8Rg",
      "type":"PsmsBlsSignature2022",
      "verificationMethod":"did:key:zAbekb81CvHVLAzkAq7C1GvGTibogjUSHeUqL7qtfd18Hh9D8DGggVj1YJiy4HAAnawQYh1FvVLVLjZ1vJ3DJCDrhzHUAaENUnspm47k9qH5jFjh9Tq4oMQze5HGpofc1orbT1A7H4VWkPBgLKPHJDG2x3T1G62NYP2MdyCUaRJ249Z8dPLyVSfHzWsiAeLR55M4sQCUxGkZUYAXo7XszbVEgHCyJCnDY3FTgmzSncbYihH845y6c85a8A2pitFGK8xmHLZttJpsBL4Er1mrSmrDN6GNmy3tnG4rpuixTGXYWGEQkERNQfNo5qdiaGbDzBR28vRup9aEgqMsbJHugEFjD1u1j9P9RxNdaASxwrDL2rsppzYKV2q3V9qFhGX5nnodff4xSZ5hoo9FvSvN2jfJvz7oEn1mx46VGZvbGPf4Kvqk4EQbZcYtwUg8ucY4871tVYyqS1Zh4QBrdrKqwVR7EnhmvNYvgZLFzWkQH6pifBfGMcUmY4o8nYcEGfdgYdY54gRh75Gfup3EobnA5FBSfMXCZMkLw2wLa7uWBkJeXRHbS3rNYaaR3FBXExdcRJARGWKWTvbRRqgYeQgW2MgZHRjqvfJuFifk88jKhiPHJMRXKBEYrXMR9LuXbbKkTDxM9JVdYxHkDTMSGPGuxnj4GXnjvfAvdh6A4x2Wb8pComHMCHdYTiQnxWzMEm3c7D8bRGPReuGoh85rSDorZ62S87ogtSuQtyQeyYQ8StcUTwpuLE6XwDbgvKNxfFroy96tQVRArEHNjc61oGup3iZom6RutczdbnzTzKopiHR2F2dsMKJD5rfQ9dSmUF3JA61ME5d6euN8DtN365rcRGP1nVnyiPUejMsLxcBBduuqCSNUnzKZ6Jfx8ZE4ssCB5rSD1t3WmPLJvQTkBagbyiHBLPveihvQ9peVXX2xzr125h1j2a4S5VUZotfnWhiyX5cVno1zNLBS8xbjB7hvQYpBe2zEtdX35xannnGdhTVaq15SZBo1sbBxx1nDYTfi7vLd2L6aBHLj4dhwPASJuFkDKo3XeRJw8nt59VhQXoLC#zAbekb81CvHVLAzkAq7C1GvGTibogjUSHeUqL7qtfd18Hh9D8DGggVj1YJiy4HAAnawQYh1FvVLVLjZ1vJ3DJCDrhzHUAaENUnspm47k9qH5jFjh9Tq4oMQze5HGpofc1orbT1A7H4VWkPBgLKPHJDG2x3T1G62NYP2MdyCUaRJ249Z8dPLyVSfHzWsiAeLR55M4sQCUxGkZUYAXo7XszbVEgHCyJCnDY3FTgmzSncbYihH845y6c85a8A2pitFGK8xmHLZttJpsBL4Er1mrSmrDN6GNmy3tnG4rpuixTGXYWGEQkERNQfNo5qdiaGbDzBR28vRup9aEgqMsbJHugEFjD1u1j9P9RxNdaASxwrDL2rsppzYKV2q3V9qFhGX5nnodff4xSZ5hoo9FvSvN2jfJvz7oEn1mx46VGZvbGPf4Kvqk4EQbZcYtwUg8ucY4871tVYyqS1Zh4QBrdrKqwVR7EnhmvNYvgZLFzWkQH6pifBfGMcUmY4o8nYcEGfdgYdY54gRh75Gfup3EobnA5FBSfMXCZMkLw2wLa7uWBkJeXRHbS3rNYaaR3FBXExdcRJARGWKWTvbRRqgYeQgW2MgZHRjqvfJuFifk88jKhiPHJMRXKBEYrXMR9LuXbbKkTDxM9JVdYxHkDTMSGPGuxnj4GXnjvfAvdh6A4x2Wb8pComHMCHdYTiQnxWzMEm3c7D8bRGPReuGoh85rSDorZ62S87ogtSuQtyQeyYQ8StcUTwpuLE6XwDbgvKNxfFroy96tQVRArEHNjc61oGup3iZom6RutczdbnzTzKopiHR2F2dsMKJD5rfQ9dSmUF3JA61ME5d6euN8DtN365rcRGP1nVnyiPUejMsLxcBBduuqCSNUnzKZ6Jfx8ZE4ssCB5rSD1t3WmPLJvQTkBagbyiHBLPveihvQ9peVXX2xzr125h1j2a4S5VUZotfnWhiyX5cVno1zNLBS8xbjB7hvQYpBe2zEtdX35xannnGdhTVaq15SZBo1sbBxx1nDYTfi7vLd2L6aBHLj4dhwPASJuFkDKo3XeRJw8nt59VhQXoLC"
   },
   "type":[
      "VerifiableCredential",
      "NancyCredential"
   ]
    }
}'
```
DeriveCredential:
```
curl --location 'http://localhost:8082/nancy/idm/deriveProof' \
--header 'Content-Type: application/json' \
--data-raw '{
    "credId": "http://example/credentials/18723",
    "nonce": "pavo",
    "querybyframe": {
        "@context": [
        "https://www.w3.org/2018/credentials/v1",
        "https://ssiproject.inf.um.es/security/psms/v1",
        "https://nancy-identity/context/exampleContext/v1",
        "https://w3id.org/security/bbs/v1"
        ],
        "type": ["VerifiableCredential", "NancyCredential"],
        "@explicit": true,
        "identifier": {},
        "issuer": {},
        "issuanceDate": {},
        "credentialSubject": {
        "@explicit": true,
        "nancyId": {},
        "acces5g": {}
        }
    }
}
'
```
VerifyCredential:
```
curl --location 'http://localhost:8082/nancy/idm/verifyCredential' \
--header 'Content-Type: application/json' \
--data-raw '{
   "credential":{
        "@context": [
            "https://www.w3.org/2018/credentials/v1",
            "https://ssiproject.inf.um.es/security/psms/v1",
            "https://nancy-identity/context/exampleContext/v1",
            "https://w3id.org/security/bbs/v1"
        ],
        "credentialSubject": {
            "acces5g": [
                "99940",
                "99941"
            ],
            "id": "did:example:ebfeb1f712ebc6f1c276e12ec21",
            "nancyId": "zbASDASDjfasfasfasfiasfjasfjioasfjas0asjfebfeb1f712ebc6f1c276e12ec21"
        },
        "id": "http://example/credentials/18723",
        "issuanceDate": "2023-07-08T19:23:24Z",
        "issuer": "did:erat:example",
        "proof": {
            "created": "2025-02-24T14:11:53.741729213+01:00",
            "nonce": "cGF2bw==",
            "proofPurpose": "assertionMethod",
            "proofValue": "AAMjAwQGD9SB96vTXzxKygwbjAYvOczSEZezqLzx06/AJGupFC/JO2ZlAFyqDeExg6OfKEUP83YeDxTihu862U3W3AgBauY9PDfwf1iRQhNGSqd+t9HJwPaKo/QrBaeXZngwF4kT6YHf7qF6ITy+ovsqh0S3LMMBTTf4r5T4VGuIEIkE9MghInYPWyq8SxYzaUxqOxAL9sOBpMZjb4vD5HOrECTCVHVlMOgohi8PQrlLDBpfm7ISCjtEDKqAvuoH6YVcPEAEAF+L5o8FJWNB7pHyGPPVFxdQwp08RgzFf1rQxkm1n1999LNrjd/RJYwlt4F8wK2BASOcmVJfydf64HrSQz2R87jxCF579IuENsKfEGj/Mg9AJe7L0GmxzZMl/bXIInUEBAp/vGT6YuCFRkbTIB6Kh1F3xbwyYH6kK9/wFncyhwvXUyIotQ8zuMPFnBkpV0ceFl0kw55CBb166/vWgJkmdXOPSDBvjeGki1rUJ36FX/vBoR7mW+q3k77pzQaPP64+BFMPUF4QXOZfd/aWr/HIAMOqA3TUD9PXi59XHKeHYe/YzZqImrBM6rJt2GED6JU7AAAAAAAAAAAAAAAAAAAAACMh+cLAvevdPxmcG8EX3bypBzALW9SJ8rH6DXb38b7vAAAAAAAAAAAAAAAAAAAAAAGIzMZ8P4Hcp0FdnrD9NFkBNl7xEzNEkmJsh/4jZ7ZfAAAAAAAAAAAAAAAAAAAAAFYPdowKTf5prRpk2bWjaclokkef8ebcfzBUkCaWnz6fAAAAAAAAAAAAAAAAAAAAABz9k7LYNkegfcfx9ROzaMdxPGs+WcNZ9EUGKQancpjYAAAAAAAAAAAAAAAAAAAAACqPuh07EieYn7aJ586puB9UkJt7r+TY1/ExWVkQI3mK",
            "type": "PsmsBlsSignatureProof2022",
            "verificationMethod": "did:key:zAbekb81CvHVLAzkAq7C1GvGTibogjUSHeUqL7qtfd18Hh9D8DGggVj1YJiy4HAAnawQYh1FvVLVLjZ1vJ3DJCDrhzHUAaENUnspm47k9qH5jFjh9Tq4oMQze5HGpofc1orbT1A7H4VWkPBgLKPHJDG2x3T1G62NYP2MdyCUaRJ249Z8dPLyVSfHzWsiAeLR55M4sQCUxGkZUYAXo7XszbVEgHCyJCnDY3FTgmzSncbYihH845y6c85a8A2pitFGK8xmHLZttJpsBL4Er1mrSmrDN6GNmy3tnG4rpuixTGXYWGEQkERNQfNo5qdiaGbDzBR28vRup9aEgqMsbJHugEFjD1u1j9P9RxNdaASxwrDL2rsppzYKV2q3V9qFhGX5nnodff4xSZ5hoo9FvSvN2jfJvz7oEn1mx46VGZvbGPf4Kvqk4EQbZcYtwUg8ucY4871tVYyqS1Zh4QBrdrKqwVR7EnhmvNYvgZLFzWkQH6pifBfGMcUmY4o8nYcEGfdgYdY54gRh75Gfup3EobnA5FBSfMXCZMkLw2wLa7uWBkJeXRHbS3rNYaaR3FBXExdcRJARGWKWTvbRRqgYeQgW2MgZHRjqvfJuFifk88jKhiPHJMRXKBEYrXMR9LuXbbKkTDxM9JVdYxHkDTMSGPGuxnj4GXnjvfAvdh6A4x2Wb8pComHMCHdYTiQnxWzMEm3c7D8bRGPReuGoh85rSDorZ62S87ogtSuQtyQeyYQ8StcUTwpuLE6XwDbgvKNxfFroy96tQVRArEHNjc61oGup3iZom6RutczdbnzTzKopiHR2F2dsMKJD5rfQ9dSmUF3JA61ME5d6euN8DtN365rcRGP1nVnyiPUejMsLxcBBduuqCSNUnzKZ6Jfx8ZE4ssCB5rSD1t3WmPLJvQTkBagbyiHBLPveihvQ9peVXX2xzr125h1j2a4S5VUZotfnWhiyX5cVno1zNLBS8xbjB7hvQYpBe2zEtdX35xannnGdhTVaq15SZBo1sbBxx1nDYTfi7vLd2L6aBHLj4dhwPASJuFkDKo3XeRJw8nt59VhQXoLC#zAbekb81CvHVLAzkAq7C1GvGTibogjUSHeUqL7qtfd18Hh9D8DGggVj1YJiy4HAAnawQYh1FvVLVLjZ1vJ3DJCDrhzHUAaENUnspm47k9qH5jFjh9Tq4oMQze5HGpofc1orbT1A7H4VWkPBgLKPHJDG2x3T1G62NYP2MdyCUaRJ249Z8dPLyVSfHzWsiAeLR55M4sQCUxGkZUYAXo7XszbVEgHCyJCnDY3FTgmzSncbYihH845y6c85a8A2pitFGK8xmHLZttJpsBL4Er1mrSmrDN6GNmy3tnG4rpuixTGXYWGEQkERNQfNo5qdiaGbDzBR28vRup9aEgqMsbJHugEFjD1u1j9P9RxNdaASxwrDL2rsppzYKV2q3V9qFhGX5nnodff4xSZ5hoo9FvSvN2jfJvz7oEn1mx46VGZvbGPf4Kvqk4EQbZcYtwUg8ucY4871tVYyqS1Zh4QBrdrKqwVR7EnhmvNYvgZLFzWkQH6pifBfGMcUmY4o8nYcEGfdgYdY54gRh75Gfup3EobnA5FBSfMXCZMkLw2wLa7uWBkJeXRHbS3rNYaaR3FBXExdcRJARGWKWTvbRRqgYeQgW2MgZHRjqvfJuFifk88jKhiPHJMRXKBEYrXMR9LuXbbKkTDxM9JVdYxHkDTMSGPGuxnj4GXnjvfAvdh6A4x2Wb8pComHMCHdYTiQnxWzMEm3c7D8bRGPReuGoh85rSDorZ62S87ogtSuQtyQeyYQ8StcUTwpuLE6XwDbgvKNxfFroy96tQVRArEHNjc61oGup3iZom6RutczdbnzTzKopiHR2F2dsMKJD5rfQ9dSmUF3JA61ME5d6euN8DtN365rcRGP1nVnyiPUejMsLxcBBduuqCSNUnzKZ6Jfx8ZE4ssCB5rSD1t3WmPLJvQTkBagbyiHBLPveihvQ9peVXX2xzr125h1j2a4S5VUZotfnWhiyX5cVno1zNLBS8xbjB7hvQYpBe2zEtdX35xannnGdhTVaq15SZBo1sbBxx1nDYTfi7vLd2L6aBHLj4dhwPASJuFkDKo3XeRJw8nt59VhQXoLC"
        },
        "type": [
            "NancyCredential",
            "VerifiableCredential"
        ]
    }
}'
```
VerifyCredential Failure:
```
curl --location 'http://localhost:8082/nancy/idm/verifyCredential' \
--header 'Content-Type: application/json' \
--data-raw '{
   "credential":{
        "@context": [
            "https://www.w3.org/2018/credentials/v1",
            "https://ssiproject.inf.um.es/security/psms/v1",
            "https://nancy-identity/context/exampleContext/v1",
            "https://w3id.org/security/bbs/v1"
        ],
        "credentialSubject": {
            "acces5g": [
                "99940",
                "99941",
                "99942"
            ],
            "id": "did:example:ebfeb1f712ebc6f1c276e12ec21",
            "nancyId": "zbASDASDjfasfasfasfiasfjasfjioasfjas0asjfebfeb1f712ebc6f1c276e12ec21"
        },
        "id": "http://example/credentials/18723",
        "issuanceDate": "2023-07-08T19:23:24Z",
        "issuer": "did:erat:example",
        "proof": {
            "created": "2025-02-24T14:11:53.741729213+01:00",
            "nonce": "cGF2bw==",
            "proofPurpose": "assertionMethod",
            "proofValue": "AAMjAwQGD9SB96vTXzxKygwbjAYvOczSEZezqLzx06/AJGupFC/JO2ZlAFyqDeExg6OfKEUP83YeDxTihu862U3W3AgBauY9PDfwf1iRQhNGSqd+t9HJwPaKo/QrBaeXZngwF4kT6YHf7qF6ITy+ovsqh0S3LMMBTTf4r5T4VGuIEIkE9MghInYPWyq8SxYzaUxqOxAL9sOBpMZjb4vD5HOrECTCVHVlMOgohi8PQrlLDBpfm7ISCjtEDKqAvuoH6YVcPEAEAF+L5o8FJWNB7pHyGPPVFxdQwp08RgzFf1rQxkm1n1999LNrjd/RJYwlt4F8wK2BASOcmVJfydf64HrSQz2R87jxCF579IuENsKfEGj/Mg9AJe7L0GmxzZMl/bXIInUEBAp/vGT6YuCFRkbTIB6Kh1F3xbwyYH6kK9/wFncyhwvXUyIotQ8zuMPFnBkpV0ceFl0kw55CBb166/vWgJkmdXOPSDBvjeGki1rUJ36FX/vBoR7mW+q3k77pzQaPP64+BFMPUF4QXOZfd/aWr/HIAMOqA3TUD9PXi59XHKeHYe/YzZqImrBM6rJt2GED6JU7AAAAAAAAAAAAAAAAAAAAACMh+cLAvevdPxmcG8EX3bypBzALW9SJ8rH6DXb38b7vAAAAAAAAAAAAAAAAAAAAAAGIzMZ8P4Hcp0FdnrD9NFkBNl7xEzNEkmJsh/4jZ7ZfAAAAAAAAAAAAAAAAAAAAAFYPdowKTf5prRpk2bWjaclokkef8ebcfzBUkCaWnz6fAAAAAAAAAAAAAAAAAAAAABz9k7LYNkegfcfx9ROzaMdxPGs+WcNZ9EUGKQancpjYAAAAAAAAAAAAAAAAAAAAACqPuh07EieYn7aJ586puB9UkJt7r+TY1/ExWVkQI3mK",
            "type": "PsmsBlsSignatureProof2022",
            "verificationMethod": "did:key:zAbekb81CvHVLAzkAq7C1GvGTibogjUSHeUqL7qtfd18Hh9D8DGggVj1YJiy4HAAnawQYh1FvVLVLjZ1vJ3DJCDrhzHUAaENUnspm47k9qH5jFjh9Tq4oMQze5HGpofc1orbT1A7H4VWkPBgLKPHJDG2x3T1G62NYP2MdyCUaRJ249Z8dPLyVSfHzWsiAeLR55M4sQCUxGkZUYAXo7XszbVEgHCyJCnDY3FTgmzSncbYihH845y6c85a8A2pitFGK8xmHLZttJpsBL4Er1mrSmrDN6GNmy3tnG4rpuixTGXYWGEQkERNQfNo5qdiaGbDzBR28vRup9aEgqMsbJHugEFjD1u1j9P9RxNdaASxwrDL2rsppzYKV2q3V9qFhGX5nnodff4xSZ5hoo9FvSvN2jfJvz7oEn1mx46VGZvbGPf4Kvqk4EQbZcYtwUg8ucY4871tVYyqS1Zh4QBrdrKqwVR7EnhmvNYvgZLFzWkQH6pifBfGMcUmY4o8nYcEGfdgYdY54gRh75Gfup3EobnA5FBSfMXCZMkLw2wLa7uWBkJeXRHbS3rNYaaR3FBXExdcRJARGWKWTvbRRqgYeQgW2MgZHRjqvfJuFifk88jKhiPHJMRXKBEYrXMR9LuXbbKkTDxM9JVdYxHkDTMSGPGuxnj4GXnjvfAvdh6A4x2Wb8pComHMCHdYTiQnxWzMEm3c7D8bRGPReuGoh85rSDorZ62S87ogtSuQtyQeyYQ8StcUTwpuLE6XwDbgvKNxfFroy96tQVRArEHNjc61oGup3iZom6RutczdbnzTzKopiHR2F2dsMKJD5rfQ9dSmUF3JA61ME5d6euN8DtN365rcRGP1nVnyiPUejMsLxcBBduuqCSNUnzKZ6Jfx8ZE4ssCB5rSD1t3WmPLJvQTkBagbyiHBLPveihvQ9peVXX2xzr125h1j2a4S5VUZotfnWhiyX5cVno1zNLBS8xbjB7hvQYpBe2zEtdX35xannnGdhTVaq15SZBo1sbBxx1nDYTfi7vLd2L6aBHLj4dhwPASJuFkDKo3XeRJw8nt59VhQXoLC#zAbekb81CvHVLAzkAq7C1GvGTibogjUSHeUqL7qtfd18Hh9D8DGggVj1YJiy4HAAnawQYh1FvVLVLjZ1vJ3DJCDrhzHUAaENUnspm47k9qH5jFjh9Tq4oMQze5HGpofc1orbT1A7H4VWkPBgLKPHJDG2x3T1G62NYP2MdyCUaRJ249Z8dPLyVSfHzWsiAeLR55M4sQCUxGkZUYAXo7XszbVEgHCyJCnDY3FTgmzSncbYihH845y6c85a8A2pitFGK8xmHLZttJpsBL4Er1mrSmrDN6GNmy3tnG4rpuixTGXYWGEQkERNQfNo5qdiaGbDzBR28vRup9aEgqMsbJHugEFjD1u1j9P9RxNdaASxwrDL2rsppzYKV2q3V9qFhGX5nnodff4xSZ5hoo9FvSvN2jfJvz7oEn1mx46VGZvbGPf4Kvqk4EQbZcYtwUg8ucY4871tVYyqS1Zh4QBrdrKqwVR7EnhmvNYvgZLFzWkQH6pifBfGMcUmY4o8nYcEGfdgYdY54gRh75Gfup3EobnA5FBSfMXCZMkLw2wLa7uWBkJeXRHbS3rNYaaR3FBXExdcRJARGWKWTvbRRqgYeQgW2MgZHRjqvfJuFifk88jKhiPHJMRXKBEYrXMR9LuXbbKkTDxM9JVdYxHkDTMSGPGuxnj4GXnjvfAvdh6A4x2Wb8pComHMCHdYTiQnxWzMEm3c7D8bRGPReuGoh85rSDorZ62S87ogtSuQtyQeyYQ8StcUTwpuLE6XwDbgvKNxfFroy96tQVRArEHNjc61oGup3iZom6RutczdbnzTzKopiHR2F2dsMKJD5rfQ9dSmUF3JA61ME5d6euN8DtN365rcRGP1nVnyiPUejMsLxcBBduuqCSNUnzKZ6Jfx8ZE4ssCB5rSD1t3WmPLJvQTkBagbyiHBLPveihvQ9peVXX2xzr125h1j2a4S5VUZotfnWhiyX5cVno1zNLBS8xbjB7hvQYpBe2zEtdX35xannnGdhTVaq15SZBo1sbBxx1nDYTfi7vLd2L6aBHLj4dhwPASJuFkDKo3XeRJw8nt59VhQXoLC"
        },
        "type": [
            "NancyCredential",
            "VerifiableCredential"
        ]
    }
}'
```

## Debugging
Install delve: 
```
go install github.com/go-delve/delve/cmd/dlv@v1.21.2
```

Build with flags for debug instead of normally:
```
go mod tidy && GO_FLAGS="-gcflags=all=-N -l" go build -o ../../build/bin/aries-agent-rest main.go
```

Run with
```
cd <directory-aries-framework>/scripts
BIN_PATH=../build/bin bash run-5g-demo-device-debug.sh
```