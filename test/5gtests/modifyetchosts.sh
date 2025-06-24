FABRIC_IP=127.0.0.1
echo "$FABRIC_IP peer0.org1.example.com" >> /etc/hosts
echo "$FABRIC_IP peer0.org2.example.com" >> /etc/hosts
echo "$FABRIC_IP ca.org1.example.com" >> /etc/hosts
echo "$FABRIC_IP orderer.example.com" >> /etc/hosts