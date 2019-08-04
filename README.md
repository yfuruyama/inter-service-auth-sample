inter-service-auth-sample
===

Sample codes for GAE-to-GAE authentication with Identity-Aware Proxy.

## Prerequisites

* Enable Identity-Aware Proxy
* Create OAuth2 Client for Identity-Aware Proxy

## Try

1. Deploy Application

```sh
PROJECT_ID=xxx
IAP_CLIENT_ID=xxx

# IAP_CLIENT_ID is OAuth2 Client ID for Identity-Aware Proxy
sed -i -e "s/{ID_TOKEN_AUDIENCE}/$IAP_CLIENT_ID/" frontend/app.yaml

gcloud --project=$PROJECT_ID app deploy frontend/
gcloud --project=$PROJECT_ID app deploy backend/
```

2. Set IAP Policy

```sh
PROJECT_ID=xxx

# Allow public users to access the frontend service
gcloud alpha iap web add-iam-policy-binding --resource-type=app-engine --service=inter-service-auth-frontend --project=$PROJECT_ID --member=allUsers --role=roles/iap.httpsResourceAccessor

# Allow GAE apps in the same project to access the backend service
gcloud alpha iap web add-iam-policy-binding --resource-type=app-engine --service=inter-service-auth-backend --project=$PROJECT_ID --member=serviceAccount:${PROJECT_ID}@appspot.gserviceaccount.com --role=roles/iap.httpsResourceAccessor
```

3. Access Frontend

```sh
PROJECT_ID=xxx

curl https://inter-service-auth-frontend-dot-${PROJECT_ID}.appspot.com
```
