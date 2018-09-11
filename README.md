# IAP Token Generator

Generate OIDC tokens to use in HTTP requests for the `Authorization: Bearer` header to make authenticated requests to Cloud IAP-secured resources. 

Implements [this flow](https://cloud.google.com/iap/docs/authentication-howto#authenticating_from_a_service_account) outlined by the GCP documentation.

Thanks to https://github.com/b4b4r07/iap_curl for implementing the OAuth flow.

## Options

```
$ iap-token-generator -h
Generate a Bearer token for making HTTP requests to IAP-protected apps

Usage:
  iap-token-generator [flags]

Flags:
  -c, --credentials string   The service account JSON credential [GOOGLE_APPLICATION_CREDENTIALS]
  -f, --filename string      Write the token to a file
  -h, --help                 help for iap-token-generator
  -i, --id string            The IAP client ID [IAP_CLIENT_ID]
  -r, --refresh              Refresh the token on an interval
```

## Usage

1. Run as a sidecar container, refresh and output the token to a file that is read by the main application.  
1. Run with refresh and consume the token from stdout.
1. Run as a subcommand with `curl --header "Authorization: Bearer $(iap-token-generator)" ...`
