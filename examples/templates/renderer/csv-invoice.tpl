{{/* A simple HTML-based invoice layout */}}
{{ define "main" }}
<!DOCTYPE html>
<html>
  <head>
    <title>{{ if eq (index . "DOCTYPE") "INVOICE" }}Invoice{{ else }}Credit Note{{ end }}</title>
  </head>
  <body>
    <div>
      <p>Dear {{ index . "fullname" }},</p>
      <p>thanks for using our services. This {{ if eq (index . "DOCTYPE") "INVOICE" }}invoice{{ else }}credit note{{ end }}
        considers all open positions up to {{ index . "DATE" }}.</p>
      <p>The total sum amounts to <b>{{ index . "GROSS" }} $<b> including
        taxes.</p>
    </div>
  </body>
</html>
{{ end }}