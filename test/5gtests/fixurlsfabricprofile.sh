CONNECTION_PROFILE_PATH="/home/jesus/Desktop/Workspace/aries-framework-5g/test/bdd/fixtures/agent-rest/data/connection-profile.json"
sed -i 's#grpcs://localhost:7051#grpcs://peer0.org1.example.com:7051#g' $CONNECTION_PROFILE_PATH
sed -i 's#https://localhost:7054#https://ca.org1.example.com:7054#g' $CONNECTION_PROFILE_PATH