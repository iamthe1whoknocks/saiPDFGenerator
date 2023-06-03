# saiPDFGenerator

### API
#### Contract command
- request:

curl --location --request GET 'http://localhost:8080' \
&emsp;    --header 'Token: SomeToken' \
&emsp;    --header 'Content-Type: application/json' \
&emsp;    --data-raw '{"method":"getPDF","data":{"link":"https://github.com/aws/aws-sdk-go/","library":"chromedp"}}'

- response: {"link":"https://link_to_download_pdf.com"} //link to download generated pdf file
  

## Configuration

1. Edit config.yml
2. Leave common section, or change values.
3. Add any settings with any depth:
```
test: " TEST"
+ any_new_chapter:
+   any_new_paragraph:
+     any_new_config: value
```
4. To access this value in the service you can use:
```
is.Context.GetConfig("any_new_chapter.any_new_paragraph.any_new_config", "default_value").(string)
```
