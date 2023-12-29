#!/usr/bin/env jb

var s = $1`ls`.split('\n').filter(v => !v.startsWith('_') && v.endsWith('.tengo')).map(v => $1(`cat ${v}`).replaceAll('import("brook")', 'undefined')).join('\n')

var h = $1`cat _header.tengo`
var f = $1`cat _footer.tengo`
write_file('/tmp/_.tengo', `
in_brooklinks := undefined
in_dnsquery := undefined
in_address := undefined
in_httprequest := undefined
in_httpresponse := undefined
${h}
${s}
${f}
`)
$1`tengo /tmp/_.tengo`
