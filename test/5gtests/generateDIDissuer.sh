curl --location 'http://localhost:8082/nancy/idm/generateDID' \
--header 'Content-Type: application/json' \
--data-raw '{
    "keys": [
        {
            "keyType": {
                "keytype": "Ed25519VerificationKey2018"
            },
            "purpose": "Authentication"
        },
        {
            "keyType": {
                "keytype": "Bls12381G1Key2022",
                "attrs": [
                    "5"
                ]
            },
            "purpose": "AssertionMethod"
        }
    ],
    "name": "nuevo13"
}
'