package music

import (
	"fmt"
	"strings"
	"time"
)

// ChinaMobileSDC is the xxx music
type ChinaMobileSDC struct {
	Song []byte
}

// NewChinaMobileSDC returns a new ChinaMobileSDC
func NewChinaMobileSDC() *ChinaMobileSDC {
	ss := make([]string, 0)
	ss = append(ss, "POST http://sdc.10086.cn/ HTTP/1.1")
	ss = append(ss, "Host: sdc.10086.cn")
	ss = append(ss, "X-Online-Host: sdc.10086.cn")
	s := strings.Join(ss, "\r\n")
	return &ChinaMobileSDC{
		Song: []byte(s),
	}
}

// Length returns length of song
func (c *ChinaMobileSDC) Length() int {
	return len(c.Song)
}

// GetSong returns song of music
func (c *ChinaMobileSDC) GetSong() []byte {
	return c.Song
}

// GetResponse returns response when the request does not equal with the song
func (c *ChinaMobileSDC) GetResponse(request []byte) []byte {
	body := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="renderer" content="webkit">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta content="中国移动官方网站为您提供业务介绍，手机话费充值查询，套餐资费介绍及网上查询办理，号码购买，优惠购机，积分查询，优惠活动等网上自助服务。" name="Description">
    <meta content="中国移动，移动，中国移动通信，中移动，中国移动首页，手机，充话费，4G，套餐，移动iphone，4G套餐，选号" name="Keywords">
    <title>中国移动官方网站</title>
    <link id="chengeIndex" href="/css/bootstrap.index.css__index2015.css" rel="stylesheet" type="text/css" />
	<link rel="shortcut icon" href="/favicon.ico" type="images/x-icon"/>
	<!--[if lt IE 9]>
		<link rel="stylesheet" href="/css/bootstrap.index.css__index2015ie8.css" type="text/css"/>
	<![endif]-->
    <script type="text/javascript" src="/js10086/jquery-1.9.1.min.js"></script>
    <script type="text/javascript" src="/head/general_head_h5.js"></script>
	<script language="JavaScript">
		var previewProv = "bj";
	</script>
	<!--[if lt IE 9]>
		<script src="/js10086/html5.js"></script>
	<![endif]-->

</head>

<body id="bootstrap">


<!-- 页头开始 -->
<header id="head"></header>
<script type="text/javascript">renderPageHead();</script>
<!-- 页头结束 -->

<div id="main">
    <!--商城商品分类-->
    <section class="shopclass visible-lg-block">
        <h1><a href="http://shop.10086.cn/mall_100_100.html">直达移动商城</a></h1>
        <ul>
            <li><img src="/images/index/shopclass01.gif"/><a href="http://shop.10086.cn/list/101_100_100_0_0_0_0.html" target="_blank">买手机</a></li>
            <li><img src="/images/index/shopclass02.gif"/><a href="http://shop.10086.cn/list/140_100_100_0_0_0_0.html" target="_blank">办套餐</a></li>
            <li><img src="/images/index/shopclass03.gif"/><a href="http://shop.10086.cn/list/146_100_100_0_0_0_0.html" target="_blank">办业务</a></li>
            <li><img src="/images/index/shopclass05.gif"/><a href="http://shop.10086.cn/list/128_100_100_0_0_0_0.html" target="_blank">挑配件</a></li>
			<li><img src="/images/index/kdzq2.png"/><a href="http://www.10086.cn/kdzq/index/" target="_blank">家庭业务</a></li>
            <li><img src="/images/index/grzx.jpg"/><a href="http://shop.10086.cn/i/" target="_blank">个人中心</a></li>
        </ul>
    </section>


    <!--轮播图-->
    <div id="banner" class="banner">
        <!--紧急公告区-->
        <div class="alert jjgg" role="alert">
            <button type="button" class="close jjggclose" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">×</span></button>
            <a href="" target="_blank"></a>
        </div>
        <div class="bs-example" data-example-id="carousel-with-captions">
            <div id="carousel-example-captions" class="carousel slide" data-ride="carousel">
                <ol class="carousel-indicators">
                
						<li class="active" data-target="#carousel-example-captions" data-slide-to="0" title="出境用移动 流量包天安心用 "></li>
					
						<li class="" data-target="#carousel-example-captions" data-slide-to="1" title="我的流量够任性"></li>
					
						<li class="" data-target="#carousel-example-captions" data-slide-to="2" title="礼享盛惠 4G+手机节"></li>
					
						<li class="" data-target="#carousel-example-captions" data-slide-to="3" title="加量不加价"></li>
					
						<li class="" data-target="#carousel-example-captions" data-slide-to="4" title="Galaxy S8 5月18日预售开启"></li>
					
						<li class="" data-target="#carousel-example-captions" data-slide-to="5" title="中国移动光宽带"></li>
					
						<li class="" data-target="#carousel-example-captions" data-slide-to="6" title="优惠购机"></li>
					
					
                
                </ol>
                <div class="carousel-inner" role="listbox">
                
						<div class="item active"><a href="http://www.10086.cn/roaming/yewu/index/#myyw "><img data-holder-rendered="true" src="/cmpop/images/index/ad1/201705/201705091427546114ga.jpg" title="出境用移动 流量包天安心用 " alt="出境用移动 流量包天安心用 " data-src="holder.js/900x500/auto/#777:#777"></a></div>
					
						<div class="item"><a href="http://shop.10086.cn/hd/skin/custom/262032/index.html"><img data-holder-rendered="true" src="/cmpop/images/index/ad1/201704/20170428161206814SG4.jpg" title="我的流量够任性" alt="我的流量够任性" data-src="holder.js/900x500/auto/#777:#777"></a></div>
					
						<div class="item"><a href="http://shop.10086.cn/hd/skin/pro/252894/index.html"><img data-holder-rendered="true" src="/cmpop/images/index/ad1/201704/201704281426436448t6.jpg" title="礼享盛惠 4G+手机节" alt="礼享盛惠 4G+手机节" data-src="holder.js/900x500/auto/#777:#777"></a></div>
					
						<div class="item"><a href="http://service.bj.10086.cn/poffice/package/queryProductItemInfoJLB.action?PACKAGECODE=LLTYTC&amp;PRODUCTCODE=Q_JLB&amp;PACKAGEID=704&amp;PRODUCTSHOWCODE=JLB1&amp;dateFlag=&amp;isCheck"><img data-holder-rendered="true" src="/cmpop/images/index/ad1/201701/20170125170020186Fx8.jpg" title="加量不加价" alt="加量不加价" data-src="holder.js/900x500/auto/#777:#777"></a></div>
					
						<div class="item"><a href="http://shop.10086.cn/hd/skin/pro/265605/index.html "><img data-holder-rendered="true" src="/cmpop/images/index/ad1/201705/20170511171406500n4z.jpg" title="Galaxy S8 5月18日预售开启" alt="Galaxy S8 5月18日预售开å¯" data-src="holder.js/900x500/auto/#777:#777"></a></div>
					
						<div class="item"><a href="http://service.bj.10086.cn/poffice/package/showKDZQ.action?PACKAGECODE=JTKDBZ&amp;isCheck=1"><img data-holder-rendered="true" src="/cmpop/images/index/ad1/201704/20170428145123399Q4R.jpg" title="中国移动光宽带" alt="中国移动光宽带" data-src="holder.js/900x500/auto/#777:#777"></a></div>
					
						<div class="item"><a href="http://shop.10086.cn/hd/events/225512.html"><img data-holder-rendered="true" src="/cmpop/images/index/ad1/201705/20170510133314155wk3.jpg" title="优惠购机" alt="优惠购机" data-src="holder.js/900x500/auto/#777:#777"></a></div>
					
					
                
                </div>
                <a class="lbnext" href="#carousel-example-captions" role="button" data-slide="prev"></a>
                <a class="lbpre" href="#carousel-example-captions" role="button" data-slide="next"></a>
            </div>
        </div>


    </div>



    <!--快捷服务区-->
    <section class="quickbtn">
        <div class="btns">
            <ul>
                <li id="kjfw_0"><a class="q0" id="kuaijiefuwu_0" href="javascript:void(0);" btnhref="" btncontent="话费查询" btntype="0" smallIcon="/cmpop/images/index/ad1/201705/20170512170532283A7g.png" bigIcon="/cmpop/images/index/ad1/201705/20170512170605524uvZ.png">话费查询</a></li>
                <li id="kjfw_1"><a class="q1" id="kuaijiefuwu_1" target="_blank" href="javascript:void(0);" btnhref="http://service.bj.10086.cn/poffice/package/ywcx.action?PACKAGECODE=CXNEW" btncontent="流量查询" btntype="0" smallIcon="/cmpop/images/index/ad1/201705/20170512170538336d6x.png" bigIcon="/cmpop/images/index/ad1/201705/20170512170610958iHS.png">流量查询</a></li>
                <li id="kjfw_2"><a class="q2" id="kuaijiefuwu_2" target="_blank" href="javascript:void(0);" btnhref="http://www.10086.cn/roaming/index/" btncontent="国际/港澳台" btntype="0" smallIcon="/cmpop/images/index/ad1/201705/20170512170543006RR9.png" bigIcon="/cmpop/images/index/ad1/201705/20170512170616309qH6.png">国际/港澳台</a></li>
                <li id="kjfw_3"><a class="q3" id="kuaijiefuwu_3" target="_blank" href="javascript:void(0);" btnhref="http://jf.10086.cn" btncontent="积分兑换" btntype="0" smallIcon="/cmpop/images/index/ad1/201705/20170512170548236WBZ.png" bigIcon="/cmpop/images/index/ad1/201705/20170512170621521lg4.png">积分兑换</a></li>
                <li id="kjfw_4"><a class="q4" id="kuaijiefuwu_4" target="_blank" href="javascript:void(0);" btnhref="http://www.10086.cn/support/service/events/" btncontent="优惠促销" btntype="0" smallIcon="/cmpop/images/index/ad1/201705/201705121705525417Wd.png" bigIcon="/cmpop/images/index/ad1/201705/20170512170626049a9d.png">优惠促销</a></li>
                <li id="kjfw_5"><a class="q5" id="kuaijiefuwu_5" target="_blank" href="javascript:void(0);" btnhref="http://www.bj.10086.cn/service/operate/" btncontent="业务办理" btntype="1" smallIcon="/cmpop/images/index/ad1/201705/201705121709148098ZP.png" bigIcon="/cmpop/images/index/ad1/201705/20170512170918825i43.png">业务办理</a></li>
            </ul>

            <ul style="display:none;">
                <li><a class="q6" id="kuaijiefuwu_6" target="_blank" href="javascript:void(0);" btnhref="http://www.bj.10086.cn/service/fee/" btncontent="余额查询" btntype="0" smallIcon="/cmpop/images/index/ad1/201507/20150731114552175FxO.gif" bigIcon="/cmpop/images/index/ad1/201507/20150731114557005M34.gif">余额查询</a></li>
                <li><a class="q7" id="kuaijiefuwu_7" target="_blank" href="javascript:void(0);" btnhref="http://service.bj.10086.cn/poffice/package/ywcx.action?PACKAGECODE=CXNEW" btncontent="套餐余量" btntype="0" smallIcon="/cmpop/images/index/ad1/201507/20150731114610673fwU.gif" bigIcon="/cmpop/images/index/ad1/201507/20150731114616221eLQ.gif">套餐余量</a></li>
                <li><a class="q8" id="kuaijiefuwu_8" target="_blank" href="javascript:void(0);" btnhref="http://www.bj.10086.cn/service/fee/zdcx/" btncontent="账单查询" btntype="0" smallIcon="/cmpop/images/index/ad1/201507/201507311146284747wW.gif" bigIcon="/cmpop/images/index/ad1/201507/20150731114636142Zwx.gif">账单查询</a></li>
                <li><a class="q9" id="kuaijiefuwu_9" target="_blank" href="javascript:void(0);" btnhref="http://www.bj.10086.cn/service/fee/qqtxdcx/" btncontent="详单查询" btntype="0" smallIcon="/cmpop/images/index/ad1/201507/20150731114650698qNY.gif" bigIcon="/cmpop/images/index/ad1/201507/20150731114657907wnF.gif">详单查询</a></li>
                <li><a class="q10" id="kuaijiefuwu_10" target="_blank" href="javascript:void(0);" btnhref="http://service.bj.10086.cn/poffice/package/ywcx.action?PACKAGECODE=CXNEW" btncontent="已订购业务" btntype="0" smallIcon="/cmpop/images/index/ad1/201507/2015073111470740167w.gif" bigIcon="/cmpop/images/index/ad1/201507/20150731114713805Gr4.gif">已订购业务</a></li>
                <li><a class="q11" id="kuaijiefuwu_11" href="javascript:void(0);" btnhref="null" btncontent="返回" btntype="0" smallIcon="/cmpop/images/index/ad1/201507/201507311147318224MB.gif" bigIcon="/cmpop/images/index/ad1/201507/20150731114739625Mg1.gif">返回</a></li>
            </ul>
        </div>

        <form id="cz_form" method="post" target="_blank" onsubmit="javascript:fun_czsubmit2();return false;">
            <div class="chongzhi">
                <div class="czjf"><a href="http://www1.10086.cn/service/czjf/czjf.jsp" target="_blank">充值交费</a> <span></span></div>
                <div class="ins"><input data-role='none' type="hidden" id="cz_val" value="100" name="amount"><input
                        data-role='none' style="font-family:微软雅黑;" type="text" id="cz_pho" ; class="in" maxlength="11"
                        size="11" name="mobileNo" onkeyup="fun_pinfo()"
                        onfocus="if (this.value=='请输入手机号码') {this.value='';} fun_hiddenerr(); this.style.color = '#333333';"
                        onblur="if( this.value=='')this.value='请输入手机号码';" value="请输入手机号码"></div>
                <div class="insselect" id="insselect" jelist="30,50,100,300,500" jrother="null">
                </div>
                <div class="czsubmit"><input type="button" id="cz_sbtn" class="cz" onclick="fun_czsubmit2();" value="立即充值">
                    <div id="cz_div_notice" class="cz_err" onclick="fun_hiddenerr();" style="display: none">
                        请输入11位正确的移动号码
                    </div>
                </div>
            </div>
        </form>

    </section>


    <!--优惠促销-->
    <section class="yhcx">
        <!--隐藏的左右箭头-->
        <div class="yhnext" style="display: block;"></div>
        <div class="yhpre" style="display: block;"></div>
        <div class="yhcon">
            <!--滚动层-->
            <div class="yhgundong" id="yhgundong">
            
					<div class="indexbox"><a target="_blank" href="http://shop.10086.cn/i/?f=rechargecredit" class="lianjie1"><h1>话费充值9.95折</h1><h2>微信扫一扫，充值更方便</h2></a><a target="_blank" href="http://shop.10086.cn/i/?f=rechargecredit" class="lianjie2"><img title="话费充值9.95折" alt="话费充值9.95折" src="/cmpop/images/index/ad/201609/201609181350225822uf.jpg" class="indeximg"></a></div>
				
					<div class="indexbox"><a target="_blank" href="http://service.bj.10086.cn/poffice/package/queryProductItemInfoJLB.action?PACKAGECODE=LLTYTC&amp;PRODUCTCODE=Q_JLB&amp;PACKAGEID=704&amp;PRODUCTSHOWCODE=JLB1" class="lianjie1"><h1>流量加量不加价</h1><h2>月费不变，流量升级</h2></a><a target="_blank" href="http://service.bj.10086.cn/poffice/package/queryProductItemInfoJLB.action?PACKAGECODE=LLTYTC&amp;PRODUCTCODE=Q_JLB&amp;PACKAGEID=704&amp;PRODUCTSHOWCODE=JLB1" class="lianjie2"><img title="流量加量不加价" alt="流量加量不加价" src="/cmpop/images/index/ad/201611/20161118102334216C2T.jpg" class="indeximg"></a></div>
				
					<div class="indexbox"><a target="_blank" href="http://service.bj.10086.cn/poffice/package/showKDZQ.action?PACKAGECODE=JTKDBZ&amp;isCheck=1" class="lianjie1"><h1>中国移动光宽带</h1><h2>4G套餐客户宽带免费用</h2></a><a target="_blank" href="http://service.bj.10086.cn/poffice/package/showKDZQ.action?PACKAGECODE=JTKDBZ&amp;isCheck=1" class="lianjie2"><img title="中国移动光宽带" alt="中国移动光宽带" src="/cmpop/images/index/ad/201610/20161010101029564Nz4.jpg" class="indeximg"></a></div>
				
					<div class="indexbox"><a target="_blank" href="http://service.bj.10086.cn/autumn/num/commonNum/showFontPage.action?busiCode=CJRZK" class="lianjie1"><h1>4G&middot;超级日租卡</h1><h2>副卡槽神器，上网通话两不误</h2></a><a target="_blank" href="http://service.bj.10086.cn/autumn/num/commonNum/showFontPage.action?busiCode=CJRZK" class="lianjie2"><img title="4G&middot;超级日租卡" alt="4G&middot;超级日租卡" src="/cmpop/images/index/ad/201701/201701101003574098Zx.jpg" class="indeximg"></a></div>
				
					<div class="indexbox"><a target="_blank" href="http://service.bj.10086.cn/poffice/package/showMarketing.action?from=bj&amp;PACKAGECODE=yxhd_llzc&amp;isCheck=1" class="lianjie1"><h1>像充话费一样充流量</h1><h2>随充随用，月末不清零</h2></a><a target="_blank" href="http://service.bj.10086.cn/poffice/package/showMarketing.action?from=bj&amp;PACKAGECODE=yxhd_llzc&amp;isCheck=1" class="lianjie2"><img title="像充话费一样充流量" alt="像充话费一样充流量" src="/cmpop/images/index/ad/201612/20161229154808185n8u.jpg" class="indeximg"></a></div>
				
					<div class="indexbox"><a target="_blank" href="http://www.bj.10086.cn/m/khd/index.html?c=200" class="lianjie1"><h1>北京移动手机营业厅</h1><h2>充值9.95折，签到送流量</h2></a><a target="_blank" href="http://www.bj.10086.cn/m/khd/index.html?c=200" class="lianjie2"><img title="北京移动手机营业厅" alt="北京移动手机营业厅" src="/cmpop/images/index/ad/201701/201701101006179067hv.jpg" class="indeximg"></a></div>
				
					<div class="indexbox"><a target="_blank" href="http://service.bj.10086.cn/poffice/package/showpackage.action?from=bj&amp;PACKAGECODE=YKLSHYLLB730&amp;isCheck=1" class="lianjie1"><h1>9.9元享1GB流量</h1><h2>7天有效，送视频会员</h2></a><a target="_blank" href="http://service.bj.10086.cn/poffice/package/showpackage.action?from=bj&amp;PACKAGECODE=YKLSHYLLB730&amp;isCheck=1" class="lianjie2"><img title="9.9元享1GB流量" alt="9.9元享1GB流量" src="/cmpop/images/index/ad/201703/20170330171106829W7s.jpg" class="indeximg"></a></div>
				
					<div class="indexbox"><a target="_blank" href="http://service.bj.10086.cn/autumn/num/commonNum/showFontPage.action?busiCode=CJLLWK" class="lianjie1"><h1>4G&middot;超级流量王卡</h1><h2>新入网充100得200</h2></a><a target="_blank" href="http://service.bj.10086.cn/autumn/num/commonNum/showFontPage.action?busiCode=CJLLWK" class="lianjie2"><img title="4G超级流量王卡" alt="4G超级流量王卡" src="/cmpop/images/index/ad/201704/20170414144820233WO8.jpg" class="indeximg"></a></div>
            
            </div>
            <!--滚动层结束-->

        </div>

    </section>

    <!--4G专区-->
	<section class="fourg container">
	  <div class="row">
	    <div class="indextitle">
	      <span>4G专区</span>
	      <a href="http://www.10086.cn/4G/index/" target="_blank">查看更多 ></a></div>
	    <div class="col-xs-6 col-sm-4 col-md-3 col-lg-3 titlepic">
	      <a href="http://www.10086.cn/4G/plus/" target="_blank">
	        <img class="lazy" src="/images/index/300380.gif" data-original="/cmpop/images/index/ad/201605/20160517132707556e41.jpg" title="4G+ 和更佳"></a>
	    </div>
	    <section class="col-xs-6 col-sm-4 col-md-3 col-lg-3 indexboxcon">
	      <div class="indexbox" id="indexboxcon_1">
	        <a class="lianjie1" href="http://service.bj.10086.cn/poffice/jsp/portal/bjyqPage/newyear_pc_Page.jsp" target="_blank">
	          <h1>4G套餐大升级</h1>
	          <h2>话费流量二选一</h2>
	        </a>
	        <a class="lianjie2" href="http://service.bj.10086.cn/poffice/jsp/portal/bjyqPage/newyear_pc_Page.jsp" target="_blank">
	          <img class="indeximg lazy" style="right: 0px;" src="/images/index/220130.gif" data-original="/cmpop/images/index/ad/201703/20170316101447732CcB.jpg" title="4G套餐大升级"></a>
	      </div>
	      <div class="indexbox2" id="indexboxcon_4">
	        <a class="lianjie1" href="http://www.10086.cn/volte/" target="_blank">
	          <h1>高清语音（VoLTE）</h1>
	          <h2>低延时，高品质，通话上网两不误</h2>
	        </a>
	        <a class="lianjie2" href="http://www.10086.cn/volte/" target="_blank">
	          <img class="indeximg lazy" style="right: 0px;" src="/images/index/220130.gif" data-original="/cmpop/images/index/ad/201605/20160518163725923nJH.jpg" title="高清语音（VoLTE）"></a>
	      </div>
	    </section>
	    <section class="col-xs-6 col-sm-4 col-md-3 col-lg-3 indexboxcon">
	      <div class="indexbox2" id="indexboxcon_2">
	        <a class="lianjie1" href="http://www.10086.cn/4G/index/bj/#F2" target="_blank">
	          <h1>4G套餐</h1>
	          <h2>加速加量不加价</h2>
	        </a>
	        <a class="lianjie2" href="http://www.10086.cn/4G/index/bj/#F2" target="_blank">
	          <img class="indeximg lazy" style="right:0;" src="/images/index/220130.gif" data-original="/cmpop/images/index/ad/201605/201605171558380616MA.jpg" title="4G套餐"></a>
	      </div>
	      <div class="indexbox" id="indexboxcon_5">
	        <a class="lianjie1" href="http://service.bj.10086.cn/autumn/num/commonNum/showFontPage.action?busiCode=CJLLWK" target="_blank">
	          <h1>4G&middot;超级流量王卡</h1>
	          <h2>重磅来袭，充100送100</h2>
	        </a>
	        <a class="lianjie2" href="http://service.bj.10086.cn/autumn/num/commonNum/showFontPage.action?busiCode=CJLLWK" target="_blank">
	          <img class="indeximg lazy" style="right:0;" src="/images/index/220130.gif" data-original="/cmpop/images/index/ad/201704/201704181406490801I9.jpg" title="流量王卡"></a>
	      </div>
	    </section>
	    <section class="col-xs-6 col-sm-4 col-md-3 col-lg-3 indexboxcon hidden-sm">
	      <div class="indexbox" id="indexboxcon_3">
	        <a class="lianjie1" href="http://www.bj.10086.cn/service/mobile/tsjf/index.html" target="_blank">
	          <h1>提速降费</h1>
	          <h2>速度更快，资费更低，服务更佳</h2>
	        </a>
	        <a class="lianjie2" href="http://www.bj.10086.cn/service/mobile/tsjf/index.html" target="_blank">
	          <img class="indeximg lazy" style="right:0;" src="/images/index/220130.gif" data-original="/cmpop/images/index/ad/201705/20170504084913374hn3.png" title="提速降费"></a>
	      </div>
	      <div class="indexbox2 hidden-sm" id="indexboxcon_6">
	        <a class="lianjie1" href="http://www.10086.cn/4G/index/bj/#F1" target="_blank">
	          <h1>4G手机</h1>
	          <h2>4G手机　品牌、厂家手机优惠大汇</h2>
	        </a>
	        <a class="lianjie2" href="http://www.10086.cn/4G/index/bj/#F1" target="_blank">
	          <img class="indeximg lazy" style="right:0;" src="/images/index/220130.gif" data-original="/cmpop/images/index/ad/201605/20160517154449126YLK.jpg" title="4G手机"></a>
	      </div>
	    </section>
	  </div>
	</section>


    <!--买手机-->
	<section class="buyphone container">
	  <div class="row">
	    <div class="indextitle">
	      <span>买手机</span>
	      <a href="http://shop.10086.cn/list/101_100_100_0_0_0_0.html" target="_blank">查看更多 ></a></div>
	    <!--鼠标滑过indexbox的时候将indeximg图片中的样式right:0;改为right:10px;-->
	    <div class="col-xs-6 col-sm-4 col-md-3 col-lg-3 indexbox">
	      <a class="lianjie1" href="http://shop.10086.cn/goods/100_100_1033280_1023239.html" target="_blank">
	        <h2>酷派8712高配版</h2>
	        <h3>支持VoLTE</h3>
	      </a>
	      <div class="price">
	        <a class="lianjie2" href="http://shop.10086.cn/goods/100_100_1033280_1023239.html" target="_blank">
					<p><goodsId:1033280,skuId:1023239></p>
				
					 <span></span>
	        </a>
	      </div>
	      <a class="lianjie3" href="http://shop.10086.cn/goods/100_100_1033280_1023239.html" target="_blank">
	        <img class="indeximg lazy" style="right: 0px;" src="/images/index/220265.gif" data-original="/cmpop/images/index/ad/201705/20170512093557577dsW.jpg" title="酷派8712高配版"></a>
	      
			<div class="buyphonejb01" style="display:none;"></div>
	      
	    </div>
	    <div class="col-xs-6 col-sm-4 col-md-3 col-lg-3 indexbox">
	      <a class="lianjie1" href="http://shop.10086.cn/goods/100_100_1042210_1036050.html" target="_blank">
	        <h2>红米note4 高配版</h2>
	        <h3>选合约，赠话费</h3>
	      </a>
	      <div class="price">
	        <a class="lianjie2" href="http://shop.10086.cn/goods/100_100_1042210_1036050.html" target="_blank">
					<p><goodsId:1042210,skuId:1036050></p>
				
					 <span></span>
	        </a>
	      </div>
	      <a class="lianjie3" href="http://shop.10086.cn/goods/100_100_1042210_1036050.html" target="_blank">
	        <img class="indeximg lazy" style="right: 0px;" src="/images/index/220265.gif" data-original="/cmpop/images/index/ad/201703/20170330122428686soc.jpg" title="红米note4"></a>

			<div class="buyphonejb02"><a href="http://shop.10086.cn/goods/100_100_1042210_1036050.html" target="_blank">优惠</a></div>

	    </div>
	    <div class="col-xs-6 col-sm-4 col-md-3 col-lg-3 indexbox">
	      <a class="lianjie1" href="http://shop.10086.cn/goods/100_100_1042396_1036144.html" target="_blank">
	        <h2>华为P10 定制版</h2>
	        <h3>选合约，赠话费</h3>
	      </a>
	      <div class="price">
	        <a class="lianjie2" href="http://shop.10086.cn/goods/100_100_1042396_1036144.html" target="_blank">
					<p><goodsId:1042396,skuId:1036144></p>
				
					 <span></span>
	        </a>
	      </div>
	      <a class="lianjie3" href="http://shop.10086.cn/goods/100_100_1042396_1036144.html" target="_blank">
	        <img class="indeximg lazy" style="right: 0px;" src="/images/index/220265.gif" data-original="/cmpop/images/index/ad/201704/20170414173330499sfO.jpg" title="华为P10"></a>

			<div class="buyphonejb03"><a href="http://shop.10086.cn/goods/100_100_1042396_1036144.html" target="_blank">新品</a></div>

	    </div>
	    <div class="col-xs-6 col-sm-4 col-md-3 col-lg-3 indexbox hidden-sm">
	      <a class="lianjie1" href="http://shop.10086.cn/goods/100_100_1042408_1036182.html" target="_blank">
	        <h2>iPhone7 移动定制版 </h2>
	        <h3>选合约，赠话费</h3>
	      </a>
	      <div class="price">
	        <a class="lianjie2" href="http://shop.10086.cn/goods/100_100_1042408_1036182.html" target="_blank">
					<p><goodsId:1042408,skuId:1036182></p>
				
					 <span></span>
	        </a>
	      </div>
	      <a class="lianjie3" href="http://shop.10086.cn/goods/100_100_1042408_1036182.html" target="_blank">
	        <img class="indeximg lazy" style="right: 0px;" src="/images/index/220265.gif" data-original="/cmpop/images/index/ad/201703/20170331100145785ZVa.jpg" title="iPhone7"></a>
	      
			<div class="buyphonejb04"><a href="http://shop.10086.cn/goods/100_100_1042408_1036182.html" target="_blank">优惠</a></div>
	      
	    </div>
	  </div>
	</section>


    <!--业务推荐-->
	<section class="ywtj container">
	  <div class="row">
	    <div class="indextitle">
	      <span>业务推荐</span>
	      <a href="http://www.bj.10086.cn/service/operate/" target="_blank">查看更多 ></a></div>
	    <div class="col-xs-6 col-sm-4 col-md-3 col-lg-3 titlepic">
	      <a id="titlepic_ywtj" href="http://service.bj.10086.cn/autumn/num/commonNum/showFontPage.action?busiCode=CJRZK" target="_blank">
	        <img id="titlepic_ywtj_img" class=" lazy" src="/images/index/300380.gif" data-original="/cmpop/images/index/ad/201704/20170418174316597t8V.png" title="4G&middot;超级日租卡"></a>
	    </div>
	    <!--鼠标滑过indexbox的时候将图片中的样式right:0;改为right:10px;-->
	    <section class="col-xs-6 col-sm-4 col-md-3 col-lg-3 indexboxcon">
	      <div class="indexbox" id="ywtj_1">
	        <a class="lianjie1" href="http://service.bj.10086.cn/poffice/package/showpackage.action?from=bj&amp;PACKAGECODE=YKLSHYLLB730&amp;isCheck=1" target="_blank">
	          <h1>9.9元享1GB流量</h1>
	          <h2>7天有效，送视频会员</h2>
	        </a>
	        <a class="lianjie2" href="http://service.bj.10086.cn/poffice/package/showpackage.action?from=bj&amp;PACKAGECODE=YKLSHYLLB730&amp;isCheck=1" target="_blank">
	          <img class="indeximg lazy" style="right: 0px;" src="/images/index/220130.gif" data-original="/cmpop/images/index/ad/201703/20170330171707709Ev2.jpg" title="9.9元享1GB流量"></a>
	      </div>
	      <div class="indexbox" id="ywtj_4">
	        <a class="lianjie1" href="http://service.bj.10086.cn/poffice/package/showpackage.action?from=bj&amp;PACKAGECODE=HJT" target="_blank">
	          <h1>和家庭分享</h1>
	          <h2>1人定套餐，全家用4G</h2>
	        </a>
	        <a class="lianjie2" href="http://service.bj.10086.cn/poffice/package/showpackage.action?from=bj&amp;PACKAGECODE=HJT" target="_blank">
	          <img class="indeximg lazy" style="right: 0px;" src="/images/index/220130.gif" data-original="/cmpop/images/index/ad/201608/20160815153742436naO.jpg" title="和家庭分享"></a>
	      </div>
	    </section>
	    <section class="col-xs-6 col-sm-4 col-md-3 col-lg-3 indexboxcon">
	      <div class="indexbox" id="ywtj_2">
	        <a class="lianjie1" href="http://service.bj.10086.cn/poffice/package/showpackage.action?from=bj&amp;PACKAGECODE=4GFXTC&amp;productShowCode=fxtcsb" target="_blank">
	          <h1>4G飞享套餐升级版</h1>
	          <h2>长途市话一个价，18元/月起</h2>
	        </a>
	        <a class="lianjie2" href="http://service.bj.10086.cn/poffice/package/showpackage.action?from=bj&amp;PACKAGECODE=4GFXTC&amp;productShowCode=fxtcsb" target="_blank">
	          <img class="indeximg lazy" style="right:0;" src="/images/index/220130.gif" data-original="/cmpop/images/index/ad/201608/20160809181105775uug.png" title="4G飞享套餐升级版"></a>
	      </div>
	      <div class="indexbox" id="ywtj_5">
	        <a class="lianjie1" href="http://service.bj.10086.cn/poffice/package/showpackage.action?PACKAGECODE=GPRSYW&amp;productShowCode=sjllkxb" target="_blank">
	          <h1>流量可选包</h1>
	          <h2>11款超低资费流量可选包</h2>
	        </a>
	        <a class="lianjie2" href="http://service.bj.10086.cn/poffice/package/showpackage.action?PACKAGECODE=GPRSYW&amp;productShowCode=sjllkxb" target="_blank">
	          <img class="indeximg lazy" style="right:0;" src="/images/index/220130.gif" data-original="/cmpop/images/index/ad/201508/20150803103436241E2q.jpg" title="流量可选包"></a>
	      </div>
	    </section>
	    <section class="col-xs-6 col-sm-4 col-md-3 col-lg-3 indexboxcon hidden-sm">
	      <div class="indexbox" id="ywtj_3">
	        <a class="lianjie1" href="http://www.bj.10086.cn/service/myzq/index.html" target="_blank">
	          <h1>国际/港澳台漫游专区</h1>
	          <h2>开通功能/套餐，查询资费/活动</h2>
	        </a>
	        <a class="lianjie2" href="http://www.bj.10086.cn/service/myzq/index.html" target="_blank">
	          <img class="indeximg lazy" style="right:0;" src="/images/index/220130.gif" data-original="/cmpop/images/index/ad/201508/201508031031274218cc.jpg" title="国际/港澳台漫游"></a>
	      </div>
	      <div class="indexbox" id="ywtj_6">
	        <a class="lianjie1" href="http://service.bj.10086.cn/poffice/package/queryProductItemInfoJLB.action?PACKAGECODE=LLTYTC&amp;PRODUCTCODE=Q_JLB&amp;PACKAGEID=704&amp;PRODUCTSHOWCODE=JLB1&amp;dateFlag=&amp;isCheck" target="_blank">
	          <h1>流量加量不加价</h1>
	          <h2>月费不变，流量升级</h2>
	        </a>
	        <a class="lianjie2" href="http://service.bj.10086.cn/poffice/package/queryProductItemInfoJLB.action?PACKAGECODE=LLTYTC&amp;PRODUCTCODE=Q_JLB&amp;PACKAGEID=704&amp;PRODUCTSHOWCODE=JLB1&amp;dateFlag=&amp;isCheck" target="_blank">
	          <img class="indeximg lazy" style="right:0;" src="/images/index/220130.gif" data-original="/cmpop/images/index/ad/201606/20160612154321418aUY.jpg" title="流量加量不加价"></a>
	      </div>
	    </section>
	  </div>
	</section>


    <!--公告-->
	 <section class="indexgg container">
	  <span>公告：</span>
	  <ul class="row">
	    <!--pad手机上显示一条-->
	    <li class="col-xs-12 col-sm-12 col-md-6 col-lg-6" style="display:block;">
	      <a class="noticeA" href="http://www.10086.cn/aboutus/news/GroupNews/201705/t20170508_63675.htm" title="中国移动参与&ldquo;一带一路&rdquo;共建情况" target="_blank">中国移动参与&ldquo;一带一路&rdquo;共建情况</a>
	      <span class="hidden-xs">2017-05-08</span></li>
	    <li class="col-xs-12 col-sm-12 col-md-6 col-lg-6" style="display:block;">
	      <a class="noticeA" href="http://www.10086.cn/aboutus/news/pannounce/bj/201705/t20170512_63707.htm" title="关于国际/港澳台长途漫游功能升级的公告" target="_blank">关于国际/港澳台长途漫游功能升级的公告</a>
	      <span class="hidden-xs">2017-05-12</span>
	    </li>
	    <li class="col-xs-12 col-sm-12 col-md-6 col-lg-6" style="display:block;">
	      <a class="noticeA" href="http://www.10086.cn/aboutus/news/announced/201705/t20170511_63701.htm" title="关于5月11日至12日网站系统升级公告" target="_blank">关于5月11日至12日网站系统升级公告</a>
	      <span class="hidden-xs">2017-05-11</span>
	    </li>
	    <li class="col-xs-12 col-sm-12 col-md-6 col-lg-6" style="display:block;">
	      <a class="noticeA" href="http://www.10086.cn/aboutus/news/pannounce/bj/201705/t20170511_63704.htm" title="关于5月11日系统升级的公告" target="_blank">关于5月11日系统升级的公告</a>
	      <span class="hidden-xs">2017-05-11</span>
	    </li>
	  </ul>
	  <div class="qhbtn">
	    <a class="left">
	      < </a>
	        <a class="right">></a></div>
	</section>


    <!--辅助需求导航区-->
	<section class="index-help container">
	  <div class="row">
	    <div class="col-xs-12 col-sm-12 col-md-5ths col-lg-5ths help_list">
	      <h2 hotimg="null">
	        <span class="col-xs-10 col-sm-10 col-md-12 col-lg-12">服务渠道</span>
	        <!--当linkgroup显示的时候，将a中的+换为x-->
	        <a class="col-xs-2 col-sm-2 visible-xs-block visible-sm-block" href="#">+</a></h2>
	      <div class="linkgroup">
	        <a href="http://www.10086.cn/support/service/channel/entity/" target="_blank" hotimg="null">自助终端</a>
	        <a href="http://www.10086.cn/support/service/channel/10086/" target="_blank" hotimg="null">10086热线</a>
	        <a href="http://www.10086.cn/support/service/channel/entity/" target="_blank" hotimg="null">实体营业厅</a>
	        <a href="http://www.10086.cn/support/service/channel/sms/" target="_blank" hotimg="null">短信营业厅</a>
	        <a href="http://www.10086.cn/support/service/channel/online/" target="_blank" hotimg="null">网上营业厅</a>
	        <a href="http://www.10086.cn/support/service/channel/mobile/" target="_blank" hotimg="null">掌上/ 手机营业厅</a>
	      </div>
	    </div>
	    <div class="col-xs-12 col-sm-12 col-md-5ths col-lg-5ths help_list">
	      <h2 hotimg="null">
	        <span class="col-xs-10 col-sm-10 col-md-12 col-lg-12">站点导航</span>
	        <a class="col-xs-2 col-sm-2  visible-xs-block visible-sm-block" href="#">+</a></h2>
	      <div class="linkgroup">
	        <a href="http://mm.10086.cn/" target="_blank" hotimg="null">MM</a>
	        <a href="http://feixin.10086.cn/" target="_blank" hotimg="null">飞信</a>
	        <a href="http://cmpay.10086.cn/" target="_blank" hotimg="null">和包</a>
	        <a href="http://read.10086.cn/u/index" target="_blank" hotimg="null">咪咕阅读</a>
	        <a href="http://mail.10086.cn/" target="_blank" hotimg="null">139邮箱</a>
	        <a href="http://music.10086.cn/" target="_blank" hotimg="null">咪咕音乐</a>
	      </div>
	    </div>
	    <div class="col-xs-12 col-sm-12 col-md-5ths col-lg-5ths help_list">
	      <h2 hotimg="null">
	        <span class="col-xs-10 col-sm-10 col-md-12 col-lg-12">快捷服务</span>
	        <a class="col-xs-2 col-sm-2 col-md-2  visible-xs-block visible-sm-block" href="#">+</a></h2>
	      <div class="linkgroup">
	        <a href="http://service.bj.10086.cn/poffice/package/showpackage.action?PACKAGECODE=XHZPP" target="_blank" hotimg="null">携号转品牌</a>
	        <a href="http://service.bj.10086.cn/poffice/package/ywcx.action?PACKAGECODE=CXNEW" target="_blank" hotimg="null">套餐余量查询</a>
	        <a href="http://service.bj.10086.cn/poffice/package/showpackagehddzcx.action?PACKAGECODE=HDDZXXCXYW" target="_blank" hotimg="null">电子卡券查询</a>
	        <a href="http://www.bj.10086.cn/index/wenjuan/260.html" target="_blank" hotimg="null">用户满意度调研</a>
	        <a href="http://www.10086.cn/roaming/" target="_blank" hotimg="null">国际/港澳台业务</a>
	        <a href="http://service.bj.10086.cn/poffice/package/showpackage.action?PACKAGECODE=KFMMYW" target="_blank" hotimg="null">客户服务密码业务</a>
	      </div>
	    </div>
	    <div class="col-xs-12 col-sm-12 col-md-5ths col-lg-5ths help_list">
	      <h2 hotimg="null">
	        <span class="col-xs-10 col-sm-10 col-md-12 col-lg-12">产品推荐</span>
	        <a class="col-xs-2 col-sm-2  visible-xs-block visible-sm-block" href="#">+</a></h2>
	      <div class="linkgroup">
	        <a href="http://service.bj.10086.cn/poffice/package/showpackage.action?PACKAGECODE=CYYW" target="_blank" hotimg="null">彩印</a>
	        <a href="http://service.bj.10086.cn/poffice/package/ywbl.action?PACKAGECODE=4GJCYW" target="_blank" hotimg="null">和4G套餐</a>
	        <a href="http://www.bj.10086.cn/index/12530/" target="_blank" hotimg="null">彩铃下载</a>
	        <a href="http://service.bj.10086.cn/poffice/package/showpackage.action?PACKAGECODE=GQXZ" target="_blank" hotimg="null">歌曲下载</a>
	        <a href="http://service.bj.10086.cn/poffice/package/showpackage.action?PACKAGECODE=LDTX" target="_blank" hotimg="null">来电提醒</a>
	        <a href="http://service.bj.10086.cn/poffice/package/showpackagesjyd.action?PACKAGECODE=SHZS&amp;productShowCode=aaa006" target="_blank" hotimg="null">违章及时通</a>
	      </div>
	    </div>
	    <div class="col-xs-12 col-sm-12 col-md-5ths col-lg-5ths help_list">
	      <h2 hotimg="null">
	        <span class="col-xs-10 col-sm-10 col-md-12 col-lg-12">商城服务指南</span>
	        <a class="col-xs-2 col-sm-2  visible-xs-block visible-sm-block" href="#">+</a></h2>
	      <div class="linkgroup">
	        <a href="http://shop.10086.cn/help/100_100/14.html" target="_blank" hotimg="null">购物指南</a>
	        <a href="http://shop.10086.cn/help/100_100/27.html" target="_blank" hotimg="null">付款方式</a>
	        <a href="http://shop.10086.cn/help/100_100/30.html" target="_blank" hotimg="null">物流配送</a>
	        <a href="http://shop.10086.cn/help/100_100/31.html" target="_blank" hotimg="null">售后服务</a>
	        <a href="http://shop.10086.cn/help/100_100/28.html" target="_blank" hotimg="null">个人中心</a>
	        <a href="" target="_blank" hotimg=""></a>
	      </div>
	    </div>
	  </div>
	</section>

<!--问卷调查浮层-->
<!--<div class="rfu2" id="rfu2">
	<a target="_blank" href="https://sc.chinamobilesz.com/survey/7826812.do?channel=3">
        <img class="kf" src="../images/index/wjdc.png" />
    </a>
</div>-->

<!--在线客服浮层-->
<div id="rfu" class="rfu">
    <img src="/images/index/kf.gif" class="kf">

    <div class="cjwt" style="left: -20px;">
        <a target="_blank" href="http://www.bj.10086.cn/onlineservice/"><img border="0" src="/images/index/zxzx.gif"></a>
    </div>

    <div class="zxzx" style="left: -20px;">
        <a target="_blank" href="http://www.bj.10086.cn/support/businesshelp/"><img border="0" src="/images/index/cjwt.gif"></a>
    </div>

    <div class="tsjy" style="left: -20px;">
        <a target="_blank" href="http://www.10086.cn/support/selfservice/suggest/"><img border="0" src="/images/index/tsjy.gif"></a>
    </div>

</div>
</div>


<!--页脚-->
<footer id="tail"></footer>

<script type="text/javascript" src="/images/js/portal/sdc_home.js"></script>
<!--[if lt IE 9]>
	<script src="/js10086/es5-shim.js"></script>
<![endif]-->

<script type="text/javascript" src="/js10086/homepage_h5libs.js"></script>
<script type="text/javascript" src="/js10086/homepage_sttlibsv2.js"></script>

<script type="text/javascript" src="/tail/general_tail_v5.js"></script>


	
</body>
</html>

<!-- bj , 100 , 100 -->`
	ss := ""
	ss += "HTTP/1.1 200 OK\r\n"
	ss += fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123))
	ss += "Server: Apache\r\n"
	ss += fmt.Sprintf("Last-Modified: %s\r\n", time.Now().Format(time.RFC1123))
	ss += "Cache-Control: max-age=43200\r\n"
	ss += "Expires: Mon, 15 May 2117 16:53:31 GMT\r\n"
	ss += "Vary: Accept-Encoding,User-Agent\r\n"
	ss += fmt.Sprintf("Content-Length: %d\r\n", len(body))
	ss += "Content-Type: text/html\r\n"
	ss += "Connection: Keep-alive\r\n"
	ss += "Via: 1.1 ID-0002262070251166 uproxy-5\r\n"
	ss += "\r\n"
	ss += body
	return []byte(ss)
}
