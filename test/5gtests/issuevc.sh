curl --location 'http://localhost:8082/nancy/idm/acceptEnrolment' \
--header 'Content-Type: application/json' \
--data-raw '{
        "idProofs":[
            {
                "attrName":"acces5g",
                "attrValue":[
                    "99940",
                    "99941"
                ]
            },
            {
                "attrName":"accesMec",
                "attrValue":true
            },
            {
                "attrName":"accessLocEng",
                "attrValue":false
            },
            {
                "attrName":"accessV2X",
                "attrValue":true
            },
            {
                "attrName":"nancyId",
                "attrValue":"zbASDASDjfasfasfasfiasfjasfjioasfjas0asjfebfeb1f712ebc6f1c276e12ec21"
            }
            
        ]
}
'