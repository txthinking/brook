<template>
  <v-app id="inspire">
    <v-navigation-drawer
      fixed
      clipped
      app
      v-model="drawer"
    >
      <v-list dense style="padding-left:16px;">
        <template v-for="(item, i) in items">
          <v-list-group v-if="item.children" v-model="item.model" no-action>
            <v-list-tile slot="item" @click="" exact-active-class="router-link-active">
              <v-list-tile-action>
                <v-icon>{{ item.model ? item.icon : item['icon-alt'] }}</v-icon>
              </v-list-tile-action>
              <v-list-tile-content>
                <v-list-tile-title>
                  {{ item.text }}
                </v-list-tile-title>
              </v-list-tile-content>
            </v-list-tile>
            <v-list-tile
              v-for="(child, i) in item.children"
              :key="i"
              v-bind:to="child.link"
              exact
            >
              <v-list-tile-action>
                <v-icon>{{ child.icon }}</v-icon>
              </v-list-tile-action>
              <v-list-tile-content>
                <v-list-tile-title>
                  {{ child.text }}
                </v-list-tile-title>
              </v-list-tile-content>
            </v-list-tile>
          </v-list-group>
          <v-divider v-else-if="item.divider"></v-divider>
          <v-list-tile v-else v-bind:to="item.link" exact>
            <v-list-tile-action>
              <v-icon>{{ item.icon }}</v-icon>
            </v-list-tile-action>
            <v-list-tile-content>
              <v-list-tile-title>
                  {{ item.text }}
              </v-list-tile-title>
            </v-list-tile-content>
          </v-list-tile>
        </template>
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
      items: [
        { icon: 'cloud', text: 'Server', link: '/', },
        { divider: true },
        { icon: 'dehaze', text: 'Mode', link: '/mode', },
        { icon: 'palette', text: 'Builtin', link: '/builtin', },
        /*
        {
          icon: 'keyboard_arrow_up',
          'icon-alt': 'keyboard_arrow_down',
          text: 'Advanced',
          model: false,
          children: [
            { icon: '', text: 'Built-in', link: '/builtin', }
          ]
        },
        */
        { divider: true },
        { icon: 'help', text: 'Help', link: '/help', },
        { icon: 'info', text: 'About', link: '/about', },
      ]
    }),
    props: {
      source: String
    }
  }
</script>
