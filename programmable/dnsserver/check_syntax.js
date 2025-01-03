#!/usr/bin/env bun

import { $ } from 'bun'
import * as fs from 'node:fs/promises'

var s = await $`ls`.text()
var l = s.split('\n').filter(v => v.endsWith('.tengo'))
for (var i = 0; i < l.length; i++) {
    s = (await $`cat ${l[i]}`.text()).replaceAll('import("brook")', 'undefined')
    await fs.writeFile('/tmp/_.tengo', `
in_dnsservers := undefined
in_dohservers := undefined
in_dnsquery := undefined
${s}
`)
    await $`tengo /tmp/_.tengo`
}

