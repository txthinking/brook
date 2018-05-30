<template>
    <v-layout>
        <v-flex xs12 sm12 md6>
            <v-card>
                <v-container fluid>
                    <v-layout column>
                        <v-flex>
                            <v-select
                                :items="[{name:'Global', value:'global'},{name:'White List', value:'white'},{name:'Black List', value:'black'},{name:'PAC', value:'pac'},{name:'Manual', value:'manual'},]"
                                v-model.trim="o.Mode"
                                label="Mode"
                                class="input-group--focused"
                                item-text="name"
                                item-value="value"
                                ></v-select>
                        </v-flex>
                        <v-flex>
                            <v-text-field
                                label="Address"
                                placeholder="127.0.0.1:1080"
                                v-model.trim="o.Address"
                                ></v-text-field>
                        </v-flex>
                        <v-flex v-if="o.Mode == 'white' || o.Mode == 'black'">
                            <v-text-field
                                label="Domain List URL"
                                placeholder="https://... or leave it blank"
                                v-model.trim="o.DomainURL"
                                ></v-text-field>
                        </v-flex>
                        <v-flex v-if="o.Mode == 'white' || o.Mode == 'black'">
                            <v-text-field
                                label="CIDR List URL"
                                placeholder="https://... or leave it blank"
                                v-model.trim="o.CidrURL"
                                ></v-text-field>
                        </v-flex>
                        <v-flex v-if="o.Mode == 'pac'">
                            <v-text-field
                                label="PAC URL"
                                placeholder="https://... or leave it blank"
                                v-model.trim="o.PacURL"
                                ></v-text-field>
                        </v-flex>
                        <v-flex>
                            <div style="color:grey;">* DO NOT EDIT unless you know what you are doing</div>
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
            Address: '127.0.0.1:1080',
            Mode: 'pac',
            DomainURL: 'https://www.txthinking.com/pac/white.list',
            CidrURL: 'https://www.txthinking.com/pac/white_cidr.list',
            PacURL: 'https://www.txthinking.com/pac/white.pac',
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
            var s = localStorage.getItem('brook/mode');
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
            if((this.o.Mode == 'white' || this.o.Mode == 'black') && this.o.DomainURL != "" && !/http/.test(this.o.DomainURL)){
                this.girl = "Invalid Domain List URL";
                this.hey = true;
                return;
            }
            if((this.o.Mode == 'white' || this.o.Mode == 'black') && this.o.CidrURL != "" && !/http/.test(this.o.CidrURL)){
                this.girl = "Invalid CIDR List URL";
                this.hey = true;
                return;
            }
            if(this.o.Mode == 'pac' && this.o.PacURL != "" && !/http/.test(this.o.PacURL)){
                this.girl = "Invalid PAC URL";
                this.hey = true;
                return;
            }
            localStorage.setItem('brook/mode', JSON.stringify(this.o));
            this.girl = "OK";
            this.hey = true;
        },
    },
}
</script>
