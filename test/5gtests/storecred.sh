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