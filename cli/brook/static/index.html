<!DOCTYPE html>
<html>
    <head>
        <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, minimum-scale=1, user-scalable=no" />
        <title>Brook</title>
        <script async src="https://www.googletagmanager.com/gtag/js?id=G-96ENZWNBX1"></script>
        <script>
            window.dataLayer = window.dataLayer || [];
            function gtag() {
                dataLayer.push(arguments);
            }
            gtag("js", new Date());

            gtag("config", "G-96ENZWNBX1");
        </script>
        <meta charset="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="stylesheet" href="simple.min.css" />
        <script src="vue.global.prod.js"></script>
        <style>
            input {
                width: 100%;
            }
        </style>
        <script>
            window.addEventListener("DOMContentLoaded", async (e) => {
                var app = {
                    data() {
                        return {
                            link: localStorage.getItem("link") ?? "",
                            listen: localStorage.getItem("listen") ?? ":1080",
                            dnsListen: localStorage.getItem("dnsListen") ?? ":5353",
                            dnsForDefault: localStorage.getItem("dnsForDefault") ?? "8.8.8.8:53",
                            dnsForBypass: localStorage.getItem("dnsForBypass") ?? "223.5.5.5:53",
                            bypassDomainList: localStorage.getItem("bypassDomainList") ?? "",
                            bypassCIDR4List: localStorage.getItem("bypassCIDR4List") ?? "",
                            bypassCIDR6List: localStorage.getItem("bypassCIDR6List") ?? "",
                            blockDomainList: localStorage.getItem("blockDomainList") ?? "",
                            enableIPv6: localStorage.getItem("enableIPv6") ? true : false,
                            status: "disconnected",
                            ing: false,
                        };
                    },
                    async created() {
                        try {
                            var r = await fetch("/status");
                            if (r.status != 200) {
                                throw await r.text();
                            }
                            this.status = await r.text();
                        } catch (e) {
                            alert(`${e}`);
                        }
                    },
                    methods: {
                        async start() {
                            try {
                                this.ing = true;
                                var s = "";
                                if (this.link) {
                                    s += ` --link '${this.link}'`;
                                    localStorage.setItem("link", this.link);
                                } else {
                                    localStorage.setItem("link", "");
                                }
                                if (this.listen) {
                                    s += ` --listen '${this.listen}'`;
                                    localStorage.setItem("listen", this.listen);
                                } else {
                                    localStorage.setItem("listen", "");
                                }
                                if (this.dnsListen) {
                                    s += ` --dnsListen '${this.dnsListen}'`;
                                    localStorage.setItem("dnsListen", this.dnsListen);
                                } else {
                                    localStorage.setItem("dnsListen", "");
                                }
                                if (this.dnsForDefault) {
                                    s += ` --dnsForDefault '${this.dnsForDefault}'`;
                                    localStorage.setItem("dnsForDefault", this.dnsForDefault);
                                } else {
                                    localStorage.setItem("dnsForDefault", "");
                                }
                                if (this.dnsForBypass) {
                                    s += ` --dnsForBypass '${this.dnsForBypass}'`;
                                    localStorage.setItem("dnsForBypass", this.dnsForBypass);
                                } else {
                                    localStorage.setItem("dnsForBypass", "");
                                }
                                if (this.bypassDomainList) {
                                    s += ` --bypassDomainList '${this.bypassDomainList}'`;
                                    localStorage.setItem("bypassDomainList", this.bypassDomainList);
                                } else {
                                    localStorage.setItem("bypassDomainList", "");
                                }
                                if (this.bypassCIDR4List) {
                                    s += ` --bypassCIDR4List '${this.bypassCIDR4List}'`;
                                    localStorage.setItem("bypassCIDR4List", this.bypassCIDR4List);
                                } else {
                                    localStorage.setItem("bypassCIDR4List", "");
                                }
                                if (this.bypassCIDR6List) {
                                    s += ` --bypassCIDR6List '${this.bypassCIDR6List}'`;
                                    localStorage.setItem("bypassCIDR6List", this.bypassCIDR6List);
                                } else {
                                    localStorage.setItem("bypassCIDR6List", "");
                                }
                                if (this.blockDomainList) {
                                    s += ` --blockDomainList '${this.blockDomainList}'`;
                                    localStorage.setItem("blockDomainList", this.blockDomainList);
                                } else {
                                    localStorage.setItem("blockDomainList", "");
                                }
                                if (this.enableIPv6) {
                                    s += ` --enableIPv6`;
                                    localStorage.setItem("enableIPv6", "true");
                                } else {
                                    localStorage.setItem("enableIPv6", "");
                                }
                                var r = await fetch(`/start?args=${encodeURIComponent(s)}`);
                                if (r.status != 200) {
                                    throw await r.text();
                                }
                                this.status = await r.text();
                                this.ing = false;
                            } catch (e) {
                                alert(`${e}`);
                                this.ing = false;
                            }
                        },
                        async stop() {
                            try {
                                this.ing = true;
                                var r = await fetch(`/stop`);
                                if (r.status != 200) {
                                    throw await r.text();
                                }
                                this.status = await r.text();
                                this.ing = false;
                            } catch (e) {
                                alert(`${e}`);
                                this.ing = false;
                            }
                        },
                    },
                };
                Vue.createApp(app).mount("body");
            });
        </script>
    </head>
    <body>
        <header>
            <h1>Brook</h1>
            <p>brook tproxy</p>
        </header>
        <main>
            <p>
                <label>--link brook link</label><br />
                <input v-model="link" placeholder="brook://..." />
            </p>
            <p>
                <label>--listen Listen address, DO NOT contain IP</label><br />
                <input v-model="listen" placeholder=":1080" />
            </p>
            <p>
                <label>--dnsListen Start a smart DNS server</label><br />
                <input v-model="dnsListen" placeholder=":5353" />
            </p>
            <p>
                <label>--dnsForDefault DNS server for resolving domains not in bypass list</label><br />
                <input v-model="dnsForDefault" placeholder="8.8.8.8:53" />
            </p>
            <p>
                <label>--dnsForBypass DNS server for resolving domains in bypass list</label><br />
                <input v-model="dnsForBypass" placeholder="223.5.5.5:53" />
            </p>
            <p>
                <label>--bypassDomainList Suffix match mode</label><br />
                <input v-model="bypassDomainList" placeholder="/path/to/local/file/example_domain.txt" />
            </p>
            <p>
                <label>--bypassCIDR4List</label><br />
                <input v-model="bypassCIDR4List" placeholder="/path/to/local/file/example_cidr4.txt" />
            </p>
            <p>
                <label>--bypassCIDR6List</label><br />
                <input v-model="bypassCIDR6List" placeholder="/path/to/local/file/example_cidr6.txt" />
            </p>
            <p>
                <label>--blockDomainList Suffix match mode</label><br />
                <input v-model="blockDomainList" placeholder="/path/to/local/file/example_domain.txt" />
            </p>
            <p>
                <label>--enableIPv6 Your local and server must support IPv6</label><br />
                <input type="checkbox" v-model="enableIPv6" />
            </p>
            <p v-if="ing"><button disabled>Waiting...</button></p>
            <p v-if="!ing && status == 'disconnected'"><button v-on:click="start">Connect</button></p>
            <p v-if="!ing && status =='connected'"><button v-on:click="stop">Disconnect</button></p>
        </main>
        <footer>
            <p><a href="https://txthinking.com">txthinking.com</a> | <a href="https://github.com/txthinking">github.com/txthinking</a> | <a href="https://talks.txthinking.com">blog</a> | <a href="https://youtube.com/txthinking">youtube</a> | <a href="https://t.me/brookgroup">telegram</a> | <a href="https://t.me/txthinking_news">news</a></p>
        </footer>
    </body>
</html>
