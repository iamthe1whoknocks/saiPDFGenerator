# saiPDFGenerator

### API
#### Contract command
- request:

curl --location --request GET 'http://localhost:8080' \
&emsp;    --header 'Token: SomeToken' \
&emsp;    --header 'Content-Type: application/json' \
&emsp;    --data-raw '{"method":"getPDF","data":{"link":"https://github.com/aws/aws-sdk-go/","library":"chromedp"}}'

- response: {"link":"https://link_to_download_pdf.com"} 
// link to download generated pdf file
  


