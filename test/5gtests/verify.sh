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