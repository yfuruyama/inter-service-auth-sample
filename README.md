inter-service-auth-sample
===

## Deploy Application

```sh
PROJECT_ID=xxx
IAP_CLIENT_ID=xxx

sed -i -e "s/{ID_TOKEN_AUDIENCE}/$IAP_CLIENT_ID/" frontend/app.yaml

gcloud --project=$PROJECT_ID app deploy frontend/
gcloud --project=$PROJECT_ID app deploy backend-iap/
```

## Set IAP Policy

```sh
PROJECT_ID=xxx

# Allow public users to access the frontend service
gcloud alpha iap web add-iam-policy-binding --resource-type=app-engine --service=inter-service-auth-frontend --project=$PROJECT_ID --member=allUsers --role=roles/iap.httpsResourceAccessor

# Allow GAE apps in the same project to access the backend service
gcloud alpha iap web add-iam-policy-binding --resource-type=app-engine --service=inter-service-auth-backend --project=$PROJECT_ID --member=serviceAccount:${PROJECT_ID}@appspot.gserviceaccount.com --role=roles/iap.httpsResourceAccessor
```

## Access Frontend

```sh
PROJECT_ID=xxx

curl https://inter-service-auth-frontend-dot-${PROJECT_ID}.appspot.com
```
