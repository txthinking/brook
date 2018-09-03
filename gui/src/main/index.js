const {app, Notification, BrowserWindow, Tray, Menu, shell, protocol} = require('electron')
const path = require('path')
const { exec, execSync } = require('child_process')
var addrToIPPort = require('addr-to-ip-port')

if (process.env.NODE_ENV !== 'development') {
  global.__static = require('path').join(__dirname, '/static').replace(/\\/g, '\\\\')
}

let w, t
var running = false
const winURL = process.env.NODE_ENV === 'development'
  ? `http://localhost:9080`
  : `file://${__dirname}/index.html`
var brook = process.platform === 'darwin' ? 'brook_darwin_amd64' : 'brook_windows_amd64.exe'
var brook1 = process.platform === 'darwin' ? 'brook1_darwin_amd64' : 'brook1_windows_amd64.exe'
var pac = process.platform === 'darwin' ? 'pac_darwin_amd64' : 'pac_windows_amd64.exe'

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
            label: 'Brook: v20180909',
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
        try{
            execSync("chmod +x " + path.join(__static, '/' + brook))
            execSync("chmod +x " + path.join(__static, '/' + brook1))
            execSync("chmod +x " + path.join(__static, '/' + pac))
        }catch(e){
            if(Notification.isSupported()){
                (new Notification({
                    title: 'When chmod +x',
                    body: e.message,
                })).show()
            }
            return
        }
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
        try{ execSync("killall " + brook) }catch(e){}
        if (o.Mode != 'manual'){
            try{ execSync("killall " + pac) }catch(e){}
        }
    }
    if (process.platform === 'win32') {
        try{ execSync("taskkill /im "+brook+" /f /t") }catch(e){}
        if (o.Mode != 'manual'){
            try{ execSync("taskkill /im "+pac+" /f /t") }catch(e){}
        }
    }
    if (o.Mode != 'manual'){
        try{ execSync(path.join(__static, '/' + brook1) + " systemproxy -r") }catch(e){}
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
        exec(path.join(__static, '/' + brook) + " " + args.join(" "), (e, out, err)=>{
            if(e){
                if(Notification.isSupported()){
                    (new Notification({
                        title: 'Brook said',
                        body: err || e.message,
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
            exec(path.join(__static, '/' + pac) + ' ' + args.join(" "), (e, out, err)=>{
                if(e){
                    if(Notification.isSupported()){
                        (new Notification({
                            title: 'PAC server said',
                            body: err || e.message,
                        })).show()
                    }
                }
                stop(o);
            })
            pu = "http://local.txthinking.com:1980/proxy.pac";
        }
        try{
            execSync(path.join(__static, '/' + brook1) + " systemproxy -u "+pu)
        }catch(e){
            if(Notification.isSupported()){
                (new Notification({
                    title: 'PAC server said',
                    body: e.message,
                })).show()
            }
            stop(o);
        }
    }

    if (process.platform === 'darwin') {
        exec("ps -acx | grep " + brook, (e, out, err)=>{
            if(!e){
                return
            }
            runBrook()
        })
        if (o.Mode != 'manual'){
            exec("ps -acx | grep " + pac, (e, out, err)=>{
                if(!e){
                    return
                }
                runPAC()
            })
        }
    }
    if (process.platform === 'win32') {
        exec("tasklist /fo list /fi \"imagename eq "+brook+"\"", (e, out, err)=>{
            if(out.indexOf(brook) !== -1){
                return
            }
            runBrook()
        })
        if (o.Mode != 'manual'){
            exec("tasklist /fo list /fi \"imagename eq "+pac+"\"", (e, out, err)=>{
                if(out.indexOf(pac) !== -1){
                    return
                }
                runPAC()
            })
        }
    }

}
