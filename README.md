# Welcome to Minna Minny


### This is all the files for deployment
```
cp favicon.ico authz_policy.csv authz_model.conf matri.db mandm .env ~/Documents/mandm/
cp -R templates ~/Documents/mandm
cp -R static ~/Documents/mandm
```
### For android build, using termx
env GOOS=android GOARCH=arm64 go build


### For ngrok, if you want to try it 
```ngrok config add-authtoken 
cd ~/Documents/ngrok
./ngrok http --domain=kingfish-smart-formally.ngrok-free.app 8080 
```


### Find pid of port running
```
lsof -i :8080
```