<template>
    <v-layout>
        <v-flex xs12 sm12 md6>
            <v-card>
                <v-container fluid>
                    <v-layout column>
                        <v-flex>
                            <v-select
                                :items="types"
                                v-model.trim="o.Type"
                                label="Type"
                                class="input-group--focused"
                                ></v-select>
                        </v-flex>
                        <v-flex>
                            <v-text-field
                                label="Server"
                                placeholder="1.2.3.4:5"
                                v-model.trim="o.Server"
                                ></v-text-field>
                        </v-flex>
                        <v-flex>
                            <v-text-field
                                label="Password"
                                placeholder="Your server password"
                                v-model="o.Password"
                                ></v-text-field>
                        </v-flex>
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
            Type: 'Brook',
            Address: '',
            Server: '',
            Password: '',
            TCPTimeout: 60,
            TCPDeadline: 0,
            UDPDeadline: 60,
            UDPSessionTime: 60,
        },
        hey: false,
        girl: "",
        types: ["Brook", "Brook Stream", "Shadowsocks", ],
    }),

    computed: {
    },

    watch: {
    },

    created () {
        this.initialize()
        this.$http.get('https://ipapi.co/ip/')
        .withCredentials()
        .end((err, res)=>{
            switch(res.status) {
            case 200:
                this.o.Address = "[::1]:1080";
                if(res.text.indexOf(":") === -1){
                    this.o.Address = "127.0.0.1:1080";
                }
                break;
            default:
                this.girl = "Can't find IP";
                this.hey = true;
                break;
            }
        });
    },

    methods: {
        initialize () {
            var s = localStorage.getItem('brook/setting');
            if (s){
                this.o = JSON.parse(s);
            }
        },
        save () {
            if(!/.+?\:\d+/.test(this.o.Address)){
                this.girl = "Invalid Address";
                this.hey = true;
                return;
            }
            if(!/.+?\:\d+/.test(this.o.Server)){
                this.girl = "Invalid Server";
                this.hey = true;
                return;
            }
            localStorage.setItem('brook/setting', JSON.stringify(this.o));
            this.girl = "OK";
            this.hey = true;
        },
    },
}
</script>
