const {app, BrowserWindow, Tray, Menu, shell, protocol} = require('electron')
const path = require('path')
const { spawn } = require('child_process')
var addrToIPPort = require('addr-to-ip-port')

if (process.env.NODE_ENV !== 'development') {
  global.__static = require('path').join(__dirname, '/static').replace(/\\/g, '\\\\')
}

let w, t
var doing = false
var brook = null
const winURL = process.env.NODE_ENV === 'development'
  ? `http://localhost:9080`
  : `file://${__dirname}/index.html`
var exec = process.platform === 'darwin' ? 'brook_macos_amd64' : 'brook_windows_amd64.exe'

function done(){
    doing = false
}

function createWindow () {
    w = new BrowserWindow({
        width: 800,
        height: 600,
        backgroundColor: '#919191',
        center: true,
        skipTaskBar: true,
        webPreferences: {
            webSecurity: false,
        },
    })
    w.loadURL(winURL)
    //w.webContents.openDevTools()
    w.webContents.on('new-window', (e, url) =>{
          e.preventDefault();
          shell.openExternal(url);
    });
    w.on('close', (e) =>{
        if (app.quitting) {
            w = null
            return
        }
        e.preventDefault()
        w.hide()
    })
}

function createTray(){
    t = new Tray(path.join(__static, '/tray.png'))
    const contextMenu = Menu.buildFromTemplate([
        {
            label: 'Toggle',
            click: ()=>{
                if(doing){
                    return
                }
                doing = true
                w.webContents.executeJavaScript("getSetting()", false, (o)=>{
                    if(brook){
                        stop(o)
                        return
                    }
                    run(o)
                })
            },
        },
        {
            label: 'Setting',
            click: ()=>{
                w.show()
            },
        },
        {
            type: 'separator',
        },
        {
            label: 'Github',
            click: ()=>{
                shell.openExternal("https://github.com/txthinking/brook");
            },
        },
        {
            label: 'Wiki',
            click: ()=>{
                shell.openExternal("https://github.com/txthinking/brook/wiki");
            },
        },
        {
            label: 'Help',
            click: ()=>{
                shell.openExternal("https://github.com/txthinking/brook/issues");
            },
        },
        {
            label: 'Brook: v20180401',
            click: ()=>{
                shell.openExternal("https://github.com/txthinking/brook/releases");
            },
        },
        {
            type: 'separator',
        },
        {
            label: 'Quit',
            click: ()=>{
                w.webContents.executeJavaScript("getSetting()", false, (o)=>{
                    if(brook){
                        cleanAndQuit(o)
                        return
                    }
                    app.quit()
                })
            },
        },
    ])
    t.setTitle('Stopped')
    t.setToolTip('Brook: stopped')
    t.setContextMenu(contextMenu)
}

app.on('ready', ()=>{
    createWindow()
    createTray()
    w.webContents.executeJavaScript("getSetting()", false, (o)=>{
        if(o.UseWhiteTrayIcon){
            t.setImage(path.join(__static, "tray_white.png"));
        }
    })
    if (process.platform === 'darwin') {
        app.dock.hide()
        spawn("chmod", ["+x", path.join(__static, '/' + exec)])
    }
})

app.on('activate', () => {
    w.show()
})
app.on('before-quit', () => {
    app.quitting = true
})
app.on('window-all-closed', () => {
})

/**
 * Auto Updater
 *
 * Uncomment the following code below and install `electron-updater` to
 * support auto updating. Code Signing with a valid certificate is required.
 * https://simulatedgreg.gitbooks.io/electron-vue/content/en/using-electron-builder.html#auto-updating
 */

/*
import { autoUpdater } from 'electron-updater'

autoUpdater.on('update-downloaded', () => {
  autoUpdater.quitAndInstall()
})

app.on('ready', () => {
  if (process.env.NODE_ENV === 'production') autoUpdater.checkForUpdates()
})
 */

function runBrook(o){
    var client = "client"
    if (o.Type === "Brook Stream"){
        client = "streamclient"
    }
    if (o.Type === "Shadowsocks"){
        client = "ssclient"
    }
    brook = spawn(path.join(__static, '/' + exec), [
        client,
        '-l',
        o.Address,
        '-i',
        o.Address ? addrToIPPort(o.Address)[0] : "",
        '-s',
        o.Server,
        '-p',
        o.Password,
        '--tcpTimeout',
        o.TCPTimeout,
        '--tcpDeadline',
        o.TCPDeadline,
        '--udpDeadline',
        o.UDPDeadline,
        '--udpSessionTime',
        o.UDPSessionTime,
    ])
    t.setTitle('')
    t.setToolTip('Brook: started')
    brook.on('exit', (code) => {
        brook = null
        t.setTitle('Stopped')
        t.setToolTip('Brook: stopped')
    });
}

function stopBrook(){
    brook.kill()
    brook = null
    t.setTitle('Stopped')
    t.setToolTip('Brook: stopped')
}

function cleanAndQuit(o){
    if (o.AutoSystemProxy){
        var sp = spawn(path.join(__static, '/' + exec), ['systemproxy', '-r'])
        sp.on('exit', (code) => {
            if(code !== 0){
                // TODO
                return
            }
            stopBrook()
            app.quit()
            return
        });
        return
    }
    stopBrook()
    app.quit()
}

function stop(o){
    if (o.AutoSystemProxy){
        var sp = spawn(path.join(__static, '/' + exec), ['systemproxy', '-r'])
        sp.on('exit', (code) => {
            if(code !== 0){
                // TODO
                return done()
            }
            stopBrook()
            return done()
        });
        return
    }
    stopBrook()
    return done()
}

function run(o){
    if (o.AutoSystemProxy){
        var s = "SOCKS5 "+o.Address+"; SOCKS "+o.Address+"; DIRECT";
        var pac = "https://pac.txthinking.com/white/" + encodeURIComponent(s)
        if (o.UseGlobalProxyMode){
            pac = "https://pac.txthinking.com/all/" + encodeURIComponent(s)
        }
        var sp = spawn(path.join(__static, '/' + exec), ['systemproxy', '-u', pac])
        sp.on('exit', (code) => {
            if(code !== 0){
                // TODO
                return done()
            }
            runBrook(o)
            return done()
        });
        return
    }
    runBrook(o)
    return done()
}
