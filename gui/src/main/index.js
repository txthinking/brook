const {app, Notification, BrowserWindow, Tray, Menu, shell, protocol} = require('electron')
const path = require('path')
const { exec } = require('child_process')
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
        exec("chmod +x " + path.join(__static, '/' + brookcmd), (error, out, err)=>{
            if(error){
                if(Notification.isSupported()){
                    (new Notification({
                        title: 'When chmod +x brook',
                        body: err,
                    })).show()
                }
            }
        })
        exec("chmod +x " + path.join(__static, '/' + paccmd), (error, out, err)=>{
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
    if(brook){
        if (process.platform !== 'win32') {
            brook.kill();
        }
        if (process.platform === 'win32') {
            exec("taskkill /pid "+brook.pid+" /f /t", (error, out, err)=>{
                // if(error){
                //     if(Notification.isSupported()){
                //         (new Notification({
                //             title: 'When kill Brook One',
                //             body: err,
                //         })).show()
                //     }
                // }
            })
        }

        if (o.Mode != 'manual'){
            exec(path.join(__static, '/' + brookcmd) + " systemproxy -r", (error, out, err)=>{
                // if(error){
                //     if(Notification.isSupported()){
                //         (new Notification({
                //             title: 'When clean system proxy',
                //             body: err,
                //         })).show()
                //     }
                // }
            })
        }

        brook = null;
    }

    if(pac){
        if (process.platform !== 'win32') {
            pac.kill();
        }
        if (process.platform === 'win32') {
            exec("taskkill /pid "+pac.pid+" /f /t", (error, out, err)=>{
                // if(error){
                //     if(Notification.isSupported()){
                //         (new Notification({
                //             title: 'When kill PAC server',
                //             body: err,
                //         })).show()
                //     }
                // }
            })
        }
        pac = null;
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
    brook = exec(path.join(__static, '/' + brookcmd) + " " + args.join(" "), (error, out, err)=>{
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

    if (o.Mode == 'manual'){
        return done()
    }

    var pu;
    if (o.Mode == 'pac'){
        pu = o.PacURL;
    }
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
        pac = exec(path.join(__static, '/' + paccmd) + ' ' + args.join(" "), (error, out, err)=>{
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

    exec(path.join(__static, '/' + brookcmd) + " systemproxy -u "+pu, (error, out, err)=>{
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

    return done()
}
