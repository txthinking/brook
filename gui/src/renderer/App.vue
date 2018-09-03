<template>
  <v-app id="inspire">
    <v-navigation-drawer
      fixed
      clipped
      app
      v-model="drawer"
    >
      <v-list dense style="padding-left:16px;">

          <v-list-tile v-bind:to="'/'" exact>
            <v-list-tile-action>
              <v-icon>cloud</v-icon>
            </v-list-tile-action>
            <v-list-tile-content>
              <v-list-tile-title>
                  Server
              </v-list-tile-title>
            </v-list-tile-content>
          </v-list-tile>
          <v-divider></v-divider>
          <v-list-tile v-if="egg" v-bind:to="'/mode'" exact>
            <v-list-tile-action>
              <v-icon>dehaze</v-icon>
            </v-list-tile-action>
            <v-list-tile-content>
              <v-list-tile-title>
                  Mode
              </v-list-tile-title>
            </v-list-tile-content>
          </v-list-tile>
          <v-list-tile v-bind:to="'/builtin'" exact>
            <v-list-tile-action>
              <v-icon>palette</v-icon>
            </v-list-tile-action>
            <v-list-tile-content>
              <v-list-tile-title>
                  Builtin
              </v-list-tile-title>
            </v-list-tile-content>
          </v-list-tile>
          <v-divider></v-divider>
          <v-list-tile v-bind:to="'/help'" exact>
            <v-list-tile-action>
              <v-icon>help</v-icon>
            </v-list-tile-action>
            <v-list-tile-content>
              <v-list-tile-title>
                  Help
              </v-list-tile-title>
            </v-list-tile-content>
          </v-list-tile>
          <v-list-tile v-bind:to="'/about'" exact>
            <v-list-tile-action>
              <v-icon>favorite</v-icon>
            </v-list-tile-action>
            <v-list-tile-content>
              <v-list-tile-title>
                  About
              </v-list-tile-title>
            </v-list-tile-content>
          </v-list-tile>
          <v-list-tile @click="openegg" exact>
            <v-list-tile-action>
              <v-icon>info</v-icon>
            </v-list-tile-action>
            <v-list-tile-content>
              <v-list-tile-title>
                  v20180909
              </v-list-tile-title>
            </v-list-tile-content>
          </v-list-tile>

      </v-list>
    </v-navigation-drawer>
    <v-toolbar
      color="dark darken-3"
      dark
      app
      clipped-left
      fixed
    >
      <v-toolbar-title :style="$vuetify.breakpoint.smAndUp ? 'width: 300px; min-width: 250px' : 'min-width: 72px'" class="ml-0 pl-3">
        <v-toolbar-side-icon @click.stop="drawer = !drawer"></v-toolbar-side-icon>
        <span class="hidden-xs-only">Brook</span>
      </v-toolbar-title>
    </v-toolbar>
    <v-content>
      <v-container fluid fill-height>
          <keep-alive>
              <router-view/>
          </keep-alive>
      </v-container>
    </v-content>
  </v-app>
</template>

<script>
  export default {
    data: () => ({
      dialog: false,
      drawer: null,
      egg: false,
      count: 0,
    }),
    props: {
      source: String
    },

    created () {
        this.initialize()
    },

    methods: {
        initialize () {
            this.egg = localStorage.getItem('brook/egg');
        },
        openegg () {
            this.count++
            if(this.count===16){
                localStorage.setItem('brook/egg', 'Yes');
                this.egg = true;
            }
        },
    },
  }
</script>
