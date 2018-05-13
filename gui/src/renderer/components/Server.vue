<template>
    <v-layout>
        <v-flex xs12 sm12 md6>
            <v-card>
                <v-container fluid>
                    <v-layout column>
                        <v-flex>
                            <v-select
                                :items="['Brook', 'Brook Stream', 'Shadowsocks', ]"
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
                                :append-icon="passwordVisibility ? 'visibility' : 'visibility_off'"
                                :append-icon-cb="() => (passwordVisibility = !passwordVisibility)"
                                :type="passwordVisibility ? 'text' : 'password'"
                                counter
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
            Server: '',
            Password: '',
            TCPTimeout: 60,
            TCPDeadline: 0,
            UDPDeadline: 60,
            UDPSessionTime: 60,
        },
        passwordVisibility: false,
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
            var s = localStorage.getItem('brook/server');
            if (s){
                this.o = JSON.parse(s);
            }
        },
        save () {
            if(!/.+?\:\d+/.test(this.o.Server)){
                this.girl = "Invalid Server";
                this.hey = true;
                return;
            }
            localStorage.setItem('brook/server', JSON.stringify(this.o));
            this.girl = "OK";
            this.hey = true;
        },
    },
}
</script>
