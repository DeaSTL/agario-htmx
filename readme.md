# HTMX Agario


### How to run

```bash
git clone https://github.com/DeaSTL/agario-htmx
cd ./agario-htmx
go run ./
```
or 
```bash 
git clone https://github.com/DeaSTL/agario-htmx
cd ./agario-htmx
docker compose up
```

Then go into your browser and go to `localhost:8080`

if you would like to deploy this then you adjust the domain and hosted zone id in the terraform script. Then run
```bash
cd infra
terraform init
terraform plan # make sure it looks right
terraform apply
```


