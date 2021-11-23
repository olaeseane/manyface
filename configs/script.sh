ssh -i manyface-server.cer ec2-user@ec2-35-159-23-84.eu-central-1.compute.amazonaws.com
scp -i manyface-server.cer ../cmd/server/server  ec2-user@ec2-35-159-23-84.eu-central-1.compute.amazonaws.com:/home/ec2-user

ngrok http --region=eu --hostname=manyface.eu.ngrok.io 8080
ngrok tcp --region=eu --remote-addr=3.tcp.eu.ngrok.io:23623 5300

./client -ws=http://manyface.eu.ngrok.io -gs=3.tcp.eu.ngrok.io:23623 -u=user1 -p=welcome
