{{ define "main" }}
<?xml version="1.0" encoding="UTF-8"?>
<root>
    <customer>
        <number>{{ index . "CUSTNO" }}</number>
        <firstname>{{ index . "FIRSTNAME" }}</firstname>
        <lastname>{{ index . "LASTNAME" }}</lastname>
        <address>
            <street>{{ index . "STREET" }}</street>
            <zip>{{ index . "ZIP" }}</zip>
            <city>{{ index . "CITY" }}</city>
        </address>
    </customer>
{{ if eq (index . "DOCTYPE") "INVOICE" }}
    <invoice>
        <date>{{ index . "DATE" }}</date>
        <net>{{ index . "NET" }}</net>
        <gross>{{ index . "GROSS" }}</gross>
    </invoice>
{{ else }}
    <credit>
        <date>{{ index . "DATE" }}</date>
        <net>{{ index . "NET" }}</net>
        <gross>{{ index . "GROSS" }}</gross>
    </credit>
{{ end }}
</root>
{{ end }}