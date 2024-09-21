#!/usr/bin/env bun

import * as fs from 'node:fs/promises'
import { $ } from 'bun'

var f = await fs.open("../readme.md", 'w+');
await fs.write(f.fd, '# Brook\n')
await fs.write(f.fd, '<!--SIDEBAR-->\n')
await fs.write(f.fd, '<!--G-R3M673HK5V-->\n')
await fs.write(f.fd, 'A cross-platform programmable network tool.\n')
await fs.write(f.fd, '\n')
await fs.write(f.fd, '# Sponsor\n')
await fs.write(f.fd, '**❤️  [Shiliew - A network app designed for those who value their time](https://www.txthinking.com/shiliew.html)**\n')

var s = await fs.readFile('getting-started.md', { encoding: 'utf8' })
await fs.write(f.fd, s)
var s = await fs.readFile('gui.md', { encoding: 'utf8' })
await fs.write(f.fd, s)
await fs.write(f.fd, '# CLI Documentation\n')
await fs.write(f.fd, 'Each subcommand has a `--example` parameter that can print the minimal example of usage\n')
var s = await $`brook mdpage`.text()
s = s.split("\n").filter(v => !v.startsWith("[")).join("\n").replace("```\n```", "```\nbrook --help\n```").split("\n").map(v => v.startsWith("**") && !v.startsWith("**Usage") ? "- " + v : v).join('\n')
s = s.replace("### help, h", "").replace("Shows a list of commands or help for one command", "").replaceAll("- **--help, -h**: show help", "")
await fs.write(f.fd, s)
var s = await fs.readFile('example.md', { encoding: 'utf8' })
await fs.write(f.fd, s)
var s = await fs.readFile('resources.md', { encoding: 'utf8' })
await fs.write(f.fd, s)
await fs.close(f.fd)
await $`markdown ../readme.md ./index.html`
