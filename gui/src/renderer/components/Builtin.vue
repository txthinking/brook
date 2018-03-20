<template>
    <v-layout>
        <v-flex xs12 sm12 md6>
            <v-card>
                <v-container fluid>
                    <v-layout column>
                        <v-flex>
                            <v-switch
                                label="Use global proxy mode"
                                v-model="o.UseGlobalProxyMode"
                                ></v-switch>
                            <p v-if="!o.UseGlobalProxyMode">White list mode. It may take a few seconds to initialize pac after you start Brook</p>
                            <p>Only effect when enable <code>Set the system proxy automatically</code></p>
                        </v-flex>
                    </v-layout>
                </v-container>
            </v-card>
            <v-card>
                <v-container fluid>
                    <v-layout column>
                        <v-flex>
                            <v-switch
                                label="Set the system proxy automatically"
                                v-model="o.AutoSystemProxy"
                                ></v-switch>
                            <p v-if="!o.AutoSystemProxy">You can use <code>socks5://[::1]:1080</code> or <code>socks5://127.0.0.1:1080</code> by yourself.</p>
                        </v-flex>
                    </v-layout>
                </v-container>
            </v-card>
            <v-card>
                <v-container fluid>
                    <v-layout column>
                        <v-flex>
                            <v-switch
                                label="Use white icon on tray"
                                v-model="o.UseWhiteTrayIcon"
                                ></v-switch>
                            <p>Need to restart</p>
                        </v-flex>
                    </v-layout>
                </v-container>
            </v-card>
            <v-card>
                <v-container fluid>
                    <v-layout column>
                        <v-flex>
                            <v-btn class="info" @click="save">Save</v-btn>
                        </v-flex>
                    </v-layout>
                </v-container>
            </v-card>
        </v-flex>
        <v-snackbar v-model="hey">
            {{girl}}
        </v-snackbar>
    </v-layout>
</template>
<script>
export default {
    data: () => ({
        o: {
            UseGlobalProxyMode: false,
            AutoSystemProxy: true,
            UseWhiteTrayIcon: false,
        },
        hey: false,
        girl: "",
    }),

    computed: {
    },

    watch: {
    },

    created () {
        this.initialize()
    },

    methods: {
        initialize () {
            var s = localStorage.getItem('BuiltIn');
            if (s){
                this.o = JSON.parse(s);
            }
        },
        save () {
            localStorage.setItem('BuiltIn', JSON.stringify(this.o));
            this.girl = "OK";
            this.hey = true;
        },
    },
}
</script>
