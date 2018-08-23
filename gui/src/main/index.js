const {app, Notification, BrowserWindow, Tray, Menu, shell, protocol} = require('electron')
const path = require('path')
const { exec } = require('child_process')
var addrToIPPort = require('addr-to-ip-port')

if (process.env.NODE_ENV !== 'development') {
  global.__static = require('path').join(__dirname, '/static').replace(/\\/g, '\\\\')
}

let w, t
var running = false
const winURL = process.env.NODE_ENV === 'development'
  ? `http://localhost:9080`
  : `file://${__dirname}/index.html`
var brook = process.platform === 'darwin' ? 'brook_macos_amd64' : 'brook_windows_amd64.exe'
var brook1 = process.platform === 'darwin' ? 'brook1_macos_amd64' : 'brook1_windows_amd64.exe'
var pac = process.platform === 'darwin' ? 'pac_macos_amd64' : 'pac_windows_amd64.exe'

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
                w.webContents.executeJavaScript("getSetting()", false, (o)=>{
                    if(running){
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
            label: 'Brook: v20180707',
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
        exec("chmod +x " + path.join(__static, '/' + brook), (error, out, err)=>{
            if(error){
                if(Notification.isSupported()){
                    (new Notification({
                        title: 'When chmod +x',
                        body: err,
                    })).show()
                }
            }
        })
        exec("chmod +x " + path.join(__static, '/' + brook1), (error, out, err)=>{
            if(error){
                if(Notification.isSupported()){
                    (new Notification({
                        title: 'When chmod +x',
                        body: err,
                    })).show()
                }
            }
        })
        exec("chmod +x " + path.join(__static, '/' + pac), (error, out, err)=>{
            if(error){
                if(Notification.isSupported()){
                    (new Notification({
                        title: 'When chmod +x pac',
                        body: err,
                    })).show()
                }
            }
        })
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

function stop(o){
    if (process.platform === 'darwin') {
        exec("ps -acx | grep " + brook, (error, out, err)=>{
            if(error){
                return
            }
            exec("killall " + brook)
        })
        if (o.Mode != 'manual'){
            exec("ps -acx | grep " + pac, (error, out, err)=>{
                if(error){
                    return
                }
                exec("killall " + pac)
            })
        }
    }
    if (process.platform === 'win32') {
        exec("tasklist /fo list /fi \"imagename eq "+brook+"\"", (error, out, err)=>{
            if(out.indexOf(brook) === -1){
                return
            }
            exec("taskkill /im "+brook+" /f /t")
        })
        if (o.Mode != 'manual'){
            exec("tasklist /fo list /fi \"imagename eq "+pac+"\"", (error, out, err)=>{
                if(out.indexOf(pac) === -1){
                    return
                }
                exec("taskkill /im "+pac+" /f /t")
            })
        }
    }
    if (o.Mode != 'manual'){
        exec(path.join(__static, '/' + brook1) + " systemproxy -r", (error, out, err)=>{
            if(error){
                if(Notification.isSupported()){
                    (new Notification({
                        title: 'When clean system proxy',
                        body: err,
                    })).show()
                }
            }
        })
    }

    running = false
    t.setTitle('Stopped')
    t.setToolTip('Brook: stopped')
}

function run(o){
    t.setTitle('')
    t.setToolTip('Brook: started')
    running = true

    var runBrook = function(){
        var client = "client"
        if (o.Type === "Brook Stream"){
            client = "streamclient"
        }
        if (o.Type === "Shadowsocks"){
            client = "ssclient"
        }
        var args = [
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
        ]
        exec(path.join(__static, '/' + brook) + " " + args.join(" "), (error, out, err)=>{
            if(error && err){
                if(Notification.isSupported()){
                    (new Notification({
                        title: 'Brook said',
                        body: err,
                    })).show()
                }
            }
            stop(o);
        })
    }
    var runPAC = function(){
        var pu = o.PacURL;
        if (o.Mode != 'pac'){
            var p = "\"SOCKS5 "+o.Address+"; SOCKS "+o.Address+"; DIRECT\"";
            var args = [
                '-m',
                o.Mode,
                '-p',
                p,
                '-l',
                ':1980',
            ]
            if(o.Mode != 'global' && o.DomainURL != ''){
                args = args.concat(['-d', o.DomainURL])
            }
            if(o.Mode != 'global' && o.CidrURL != ''){
                args = args.concat(['-c', o.CidrURL])
            }
            exec(path.join(__static, '/' + pac) + ' ' + args.join(" "), (error, out, err)=>{
                console.log(error, err)
                if(error && err){
                    if(Notification.isSupported()){
                        (new Notification({
                            title: 'PAC server said',
                            body: err,
                        })).show()
                    }
                }
                stop(o);
            })
            pu = "http://local.txthinking.com:1980/proxy.pac";
        }
        exec(path.join(__static, '/' + brook1) + " systemproxy -u "+pu, (error, out, err)=>{
            if(error){
                if(Notification.isSupported()){
                    (new Notification({
                        title: 'When set system proxy',
                        body: err,
                    })).show()
                }
                stop(o);
            }
        })
    }

    if (process.platform === 'darwin') {
        exec("ps -acx | grep " + brook, (error, out, err)=>{
            if(!error){
                return
            }
            runBrook()
        })
        if (o.Mode != 'manual'){
            exec("ps -acx | grep " + pac, (error, out, err)=>{
                if(!error){
                    return
                }
                runPAC()
            })
        }
    }
    if (process.platform === 'win32') {
        exec("tasklist /fo list /fi \"imagename eq "+brook+"\"", (error, out, err)=>{
            if(out.indexOf(brook) !== -1){
                return
            }
            runBrook()
        })
        if (o.Mode != 'manual'){
            exec("tasklist /fo list /fi \"imagename eq "+pac+"\"", (error, out, err)=>{
                if(out.indexOf(pac) !== -1){
                    return
                }
                runPAC()
            })
        }
    }

}
