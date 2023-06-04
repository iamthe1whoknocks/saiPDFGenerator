# saiPDFGenerator

### API
- request:

curl --location --request GET 'http://localhost:8080' \
&emsp;    --header 'Token: SomeToken' \
&emsp;    --header 'Content-Type: application/json' \
&emsp;    --data-raw '{"method":"getPDF","data":{"html":"$base64 encoded byte array","library":"$lib to use to generate pdf"}}'

- response: {"link":"$link_to_download_pdf.com"} 
// link to download generated pdf file
  

- response: {"response":"$raw byte array"} 
// raw byte array

## available libraries
`wkhtmltopdf`
`chromedp`
