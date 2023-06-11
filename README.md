# PDF Generator
`Note: some html content required additional PDF converter settings`

### API
- request:

curl --location --request GET 'http://localhost:8080' \
&emsp;    --header 'Token: SomeToken' \
&emsp;    --header 'Content-Type: application/json' \
&emsp;    --data-raw '{  
&emsp;&emsp;  "method":"getPDF","data":{  
&emsp;&emsp;&emsp;    "html":"$base64_encoded_html",  
&emsp;&emsp;&emsp;    "library":"$lib"  
&emsp;&emsp;&emsp;    "output":"$outputType"  
&emsp;&emsp;  }  
&emsp;  }'  

- response: {"result":"$link_to_download_pdf", "status": "Ok"} 
// link to download generated pdf file

## available outputs: $outputType
`file`
`s3`

## available libraries: $lib
`wkhtmltopdf`
`chromedp`

## example base64: $base64_encoded_html
CjwhRE9DVFlQRSBodG1sPgo8aHRtbD4KPGhlYWQ+CiAgICAgPHRpdGxlPiBDSEVXSUUgPC90aXRsZT4KPC9oZWFkPgogICA8Ym9keSBzdHlsZT0iYmFja2dyb3VuZC1jb2xvcjojZmZlZWE1OyBmb250LWZhbWlseTp2ZXJkYW5hOyAiPgogICAgICAgIDxoMSBzdHlsZT0iY29sb3I6IzRjNmE5YjsgdGV4dC1hbGlnbjogY2VudGVyIj5DaGV3aWUncyBXZWJzaXRlPC9oMT4KICAgICAgICA8aW1nIHNyYz0icGljLmpwZyIgIHdpZHRoPSI1MDAiIGhlaWdodD0iNTAwIj4KICAgICAgICA8cD5DaGV3aWUgaXMgYSByZWQgdG95IHBvb2RsZS4gPGJyPiBIaXMgYmlydGhkYXkgaXMgb24gSnVuZSA0LCBhbmQgYXMgb2YgMjAxOSwgaGUgaXMgMyAoZG9nKSB5ZWFycyBvbGQuPC9wPgogICAgICAgIDxoMz4gQ2hld2llJ3MgZmF2b3JpdGUgdGhpbmdzIHRvIGRvOiA8L2gzPgogICAgICAgICAgICAgIDxvbD4KICAgICAgICAgICAgICAgIDxsaT5TbGVlcGluZzwvbGk+CiAgICAgICAgICAgICAgICA8bGk+R29pbmcgb24gd2Fsa3M8L2xpPgogICAgICAgICAgICAgICAgPGxpPlBsYXlpbmcgY2F0Y2g8L2xpPgogICAgICAgICAgICAgICAgPGxpPkJhcmtpbmcgYXQgc3RyYW5nZXJzPC9saT4KICAgICAgICAgICAgICA8L29sPgogICAgICAgPGEgaHJlZj0iaHR0cHM6Ly93d3cuaW5zdGFncmFtLmNvbS9kajBqb25lZ29yby8iPkZvbGxvdyBteSBvd25lcidzIGluc3RhZ3JhbSE8L2E+CiAgPC9ib2R5Pgo8L2h0bWw+

