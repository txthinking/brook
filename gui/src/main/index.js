const {app, Notification, BrowserWindow, Tray, Menu, shell, protocol} = require('electron')
const path = require('path')
const { spawn } = require('child_process')
var addrToIPPort = require('addr-to-ip-port')

if (process.env.NODE_ENV !== 'development') {
  global.__static = require('path').join(__dirname, '/static').replace(/\\/g, '\\\\')
}

let w, t
var doing = false
var brook = null
var pac = null
const winURL = process.env.NODE_ENV === 'development'
  ? `http://localhost:9080`
  : `file://${__dirname}/index.html`
var brookcmd = process.platform === 'darwin' ? 'brook_macos_amd64' : 'brook_windows_amd64.exe'
var paccmd = process.platform === 'darwin' ? 'pac_macos_amd64' : 'pac_windows_amd64.exe'

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
            label: 'Brook: v20180601',
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
                    stop(o);
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
        spawn("chmod", ["+x", path.join(__static, '/' + brookcmd)])
        spawn("chmod", ["+x", path.join(__static, '/' + paccmd)])
    }
    if (process.platform === 'win32') {
        app.setAppUserModelId("com.txthinking.brook")
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

function stop(o){
    if(pac){
        pac.kill();
        pac = null;
    }
    if(brook){
        brook.kill();
        brook = null;
    }
    if (o.Mode != 'manual'){
        var sp = spawn(path.join(__static, '/' + brookcmd), ['systemproxy', '-r'])
        sp.on('exit', (code) => {
            if(code !== 0){
                if(process.platform === 'darwin' && Notification.isSupported()){
                    (new Notification({
                        title: 'Failed',
                        body: 'Failed to clean system proxy',
                    })).show()
                }
            }
        });
    }
    t.setTitle('Stopped')
    t.setToolTip('Brook: stopped')
    return done()
}

function run(o){
    t.setTitle('')
    t.setToolTip('Brook: started')

    var client = "client"
    if (o.Type === "Brook Stream"){
        client = "streamclient"
    }
    if (o.Type === "Shadowsocks"){
        client = "ssclient"
    }
    brook = spawn(path.join(__static, '/' + brookcmd), [
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
    brook.on('exit', (code) => {
        if(Notification.isSupported()){
            (new Notification({
                title: 'Stopped',
                body: 'Brook has stopped',
            })).show()
        }
        stop(o);
    });

    if (o.Mode == 'manual'){
        return done()
    }
    var pu;
    if (o.Mode == 'pac'){
        pu = o.PacURL;
    }
    if (o.Mode != 'pac'){
        var p = "SOCKS5 "+o.Address+"; SOCKS "+o.Address+"; DIRECT";
        pac = spawn(path.join(__static, '/' + paccmd), [
            '-m',
            o.Mode,
            '-d',
            o.Mode == 'global' ? '' : o.DomainURL,
            '-c',
            o.Mode == 'global' ? '' : o.CidrURL,
            '-p',
            p,
            '-s',
            ':1980',
        ])
        pac.on('exit', (code) => {
            if(Notification.isSupported()){
                (new Notification({
                    title: 'Stopped',
                    body: 'PAC server has stopped',
                })).show()
            }
            stop(o)
        });
        pu = "http://local.txthinking.com:1980/proxy.pac";
    }
    var sp = spawn(path.join(__static, '/' + brookcmd), ['systemproxy', '-u', pu])
    sp.on('exit', (code) => {
        if(code !== 0){
            if(Notification.isSupported()){
                (new Notification({
                    title: 'Failed',
                    body: 'Failed to set system proxy',
                })).show()
            }
            stop(o)
        }
    });
    return done()
}
