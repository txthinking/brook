package music

import (
	"fmt"
	"strings"
	"time"
)

// ChinaUnicomIRead is the xxx music
type ChinaUnicomIRead struct {
	Song []byte
}

// NewChinaUnicomIRead returns a new ChinaUnicomIRead
func NewChinaUnicomIRead() *ChinaUnicomIRead {
	ss := make([]string, 0)
	ss = append(ss, "POST http://iread.wo.com.cn/ HTTP/1.1")
	ss = append(ss, "Host: iread.wo.com.cn")
	ss = append(ss, "X-Online-Host: iread.wo.com.cn")
	ss = append(ss, "Connection: keep-alive")
	ss = append(ss, "User-Agent: iread")
	ss = append(ss, "Content-Type: application/octet-stream")
	ss = append(ss, "Referer: http://iread.wo.com.cn/")
	ss = append(ss, "Accept-Encoding: gzip, deflate, br")
	ss = append(ss, "Accept-Language: zh-CN,zh;q=0.8,en-US;q=0.6,en;q=0.4")
	s := strings.Join(ss, "\r\n")
	s += "\r\n"
	return &ChinaUnicomIRead{
		Song: []byte(s),
	}
}

// Length returns length of song
func (c *ChinaUnicomIRead) Length() int {
	return len(c.Song)
}

// GetSong returns song of music
func (c *ChinaUnicomIRead) GetSong() []byte {
	return c.Song
}

// GetResponse returns response when the request does not equal with the song
func (c *ChinaUnicomIRead) GetResponse(request []byte) []byte {
	body := `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta name="renderer" content="webkit">
<meta content="text/html; charset=utf-8" http-equiv="Content-Type">
<meta name="baidu-site-verification" content="tuexnrPqBi" />
<title>中国联通沃阅读 - 小说阅读 - 免费下载 - 免流量下载 - 畅销流行小说包月</title>
<meta name="keywords" content="正版,流行,原著,畅销,大神,名家,免费下载,好看,排行,top,txt,全本,完本,连载,最新章节,更新最快">
<meta name="description" content="中国联通官方自营阅读平台,几十万册图书杂志免费免流量在线和下载阅读,包括都市,言情,穿越,玄幻,悬疑,校园,武侠,仙侠,网游各类名家大神作品,畅享阅读饕餮盛宴">
<link type="text/css" rel="stylesheet" href="/pages/2014/style/public.css">
<link type="text/css" rel="stylesheet" href="/pages/2014/style/base.css">
<link type="text/css" rel="stylesheet" href="/pages/2014/style/pages.css">
<script src="/pages/2014/js/jquery-1.11.1.min.js" type="text/javascript"></script>
<script type="text/javascript" src="/pages/2014/js/jquery.qrcode.min.js"></script>
<script src="/pages/2014/js/public.js" type="text/javascript"></script>
<script src="/pages/2014/js/js_index.js" type="text/javascript"></script>
<script src="/pages/2014/js/js_head.js" type="text/javascript" id="head" base=""></script>
<style>
	.qq-contact-cont{width:120px;height:143px;background-color:#fff;border:1px solid #ddd;padding:10px0 10px 10px;}
	.arrow-rt{position:relative;top:2px;right:-20px;z-index:1;}
	.qq-list dd{height:30px;line-height:40px;}
</style>
<!--腾讯客服-->
<script type="text/javascript">
function sideAd(){
    var ww = window.screen.availWidth;
    var hh = window.screen.availHeight;
    if (ww > 1205)
    {
	    $("#ad_right").css("visibility", "visible");
        $("#ad_right").html = '<img src="http://www.cnitpm.com/images/ad_tvyp.jpg" alt=""  width="102" height="302" />';
        $("#ad_right").css("left",ww-130);
        $("#ad_right").css("top",(hh - 702)/2 + 200);
        $(window).scroll(function() {
            var top = $(window).scrollTop();
            $("#ad_right").css("top", (hh - 702)/2 + 200 + top);
        });
    }
    else
    {
        $("#ad_right").css("visibility", "hidden");
    }
    
}

function allImageHide(imageName1,imageName2){
	$("#"+imageName1).hide();
	$("#"+imageName2).hide();
}

function showImage(showImageId,hideImageId){
	alert(showImageId);
	alert(hideImageId);
	
}
</script>

<style type="text/css">
	.index_msg{
		overflow:hidden;
		height:51px;
	}
	.index_msg_list{
		
	}
	.ran_index_flag{
		width:16px;
		height:16px;
		color:#fff;
		line-height:16px;
		font-size:12px;
		text-align:center;
		background-color:#bebebe;
		float:left;
		margin-top:8px;
		margin-right:8px;
		background-color:#b40f13;
		
	}
	
	a, a:link, a:visited, a:hover, a:active {
	    color: inherit;
	    text-decoration: none;
	}
	
	.bklist_bk{
		display:block;
		overflow:hidden;
	    text-overflow: ellipsis;
	    white-space: nowrap;
	}
	
</style>

<script>
	$(document).ready(function(){
		$("#qqkfhover").hide();
	 	$("#wechatkfhover").hide();
	 	$("#qqkfhovera").mouseover(function(){
	 		$("#qqkfhover").show();
	 		$("#wxQrcodeIDH").hide();
	 	});
		$("#wechatkfhovera").mouseover(function(){
			$("#qqkfhover").hide();
			dialogResize("#wxQrcodeIDH", true);
			$("#wxQrcodeIDH").show();
    		
		});
		$("#qqkfhovera").mouseout(function(){
	 		$("#wxQrcodeIDH").hide();
		});
		$("#wechatkfhovera").mouseout(function(){
			$("#qqkfhover").hide();
	 		$("#wxQrcodeIDH").hide();
		}); 
		
		$("#qqkfhover").mouseover(function(){
	 		$("#qqkfhover").show();
	 		$("#wxQrcodeIDH").hide();
	 	});
	 	
		$("#ad_right").mouseout(function(){
			$("#qqkfhover").hide();
	 		$("#wxQrcodeIDH").hide();
		});
		
		
		
		$("div[name='showCatindex']").hover(
			function () {
			$(this).trigger("click");
		},
		function () {
		
		});


		$("div[name='wx']").hover(
		function () {
			$(this).trigger("click");
		},
		function () {
		}
		);


		initBanner("#index-nspd_110468");	
		initBanner("#index-nspd_106635");	
		initBanner("#index-nspd_106781");	
		initBanner("#index-nspd_7880");	
		msgMove();	// xiaoxi gundong
						
	 	$.ajax({
		url : "/getReadHistory.action",
		type : 'POST',
		dataType:'html',
		timeout : 10000,
		error : function() {
			
		},
		success : function(data) {
		   $("#readAndNologin").html(data);
		   
		   var tload = false;
		   $(".index_sign_btn").click(function(){
		   		if($(this).hasClass("index_sign_signed")){
		   			return ;
		   		}
				if(tload){
					return ;
				}
				tload = true; 
				$.ajax({
					url:	  "/sign.action",
				   	cache:    false,
				   	dataType: 'json',
				   	timeout:  20000,
				   	type: 	  'POST',
				   	error: 	  function(){
				   				showMessage('网络不给力，请稍后再试');
				   			  },
				   	complete: function(){tload = false;},
				   	success:  function(result){
				   				if(result.result){	// 用户未登录或登录失效
				   					showMessage("用户未登录或登录失效",{hideBack:function(){
				   						window.location.reload(true);
				   					}});
				   				}
				   				
				   				if(result.code == '0000'){
				   					var msg = null;
				   					var giveyd = result.giveyd;
				   					if(giveyd != 0){
				   						msg = "恭喜,"+giveyd+"阅点到手啦";
				   					}else{
				   						msg = "已连续签到"+result.continuous_d+"天，累计签到"+result.accumulate_d+"天";
				   					}
				   					showMessage(msg);
				   					$(".index_sign_btn").addClass("index_sign_signed");
				   					$(".index_sign_btn span").remove();
				   					$(".index_sign_btn").text("已签到");
				   				}else{
				   					if(result.message){
				   						showMessage(result.message);
				   					}else{
				   						showMessage('网络不给力，请稍后再试...');
				   					}
				   				}
				   		}
				});
			});
			
		}
		
	});
		//增加左边菜单伸缩功能
		$(".index-l-nav").mouseover(function(e) {
			$(this).removeClass("w200");
			$("ul",$(this)).each(function(i,val){
				if(i>0){
					$(this).show();
				}
			})
			$(this).find(".main-fl-arrows").removeClass("in").addClass("out");	
			$(this).find(".icon-main-title-none").removeClass("icon-main-title-none").addClass("icon-main-title-r");
				
        });	
		$(".index-l-nav").mouseleave(function(e) {
			$(this).addClass("w200");
			$("ul",$(this)).each(function(i,val){
				if(i>0){
					$(this).hide();
				}
			})
			$(this).find(".main-fl-arrows").removeClass("out").addClass("in");	
			$(this).find(".icon-main-title-r").removeClass("icon-main-title-r").addClass("icon-main-title-none");					
        });		

   
initBrandZone();

    		//手机阅读二维码 
    		$("#wxdownloadh").mouseenter(function(){
    			dialogResize("#clientQrcodeIDH", true);
    			$("#clientQrcodeIDH").show();
    		});
    		
    		$("#wxdownloadh").mouseleave(function(){
    			$("#clientQrcodeIDH").hide();
    		});
    		
    		    		//客户端下载二维码  
    		$("#clienth").mouseenter(function(){
    			dialogResize("#wxQrcodeIDH", true);
    			$("#wxQrcodeIDH").show();
    		});
    		
    		$("#clienth").mouseleave(function(){
    			$("#wxQrcodeIDH").hide();
    		});
    		
    		
    		
    		    		
    		
	})
	function showBook(hid,obj){
	 	$("div[name='showCatindex']").attr("class","main-c-tab-cell");
		$(obj).attr("class","main-c-tab-cell slt");
	    $("div[name='catindex']").hide();
	     $('#'+hid).show();	
	}
	
	function show(hid,obj){
 $("div[name='wx']").attr("class","main-c-tab-cell");
  $(obj).attr("class","main-c-tab-cell slt");
    $("div[name='down']").hide();
     $('#'+hid).show();	
	}
		// partnerid:10023为中信出版, partnerid:10033为凤凰传媒, partnerid:10015为17K小说网, partnerid:10031为塔读文学, partnerid:10032为熊猫看书 
	  function initBrandZone(opts){	
		$.ajax({
		url : "/brandzone/getIndexBrandsInfo.action",
		type : 'POST',
		data:{ partnersid :'10005,10006,10003,10007,10004'},
		dataType:'html',
		timeout : 10000,
		error : function() {
			
		},
		success : function(data) {
		   $("#brandzone").html(data);	
$("a[name='brands']").hover(
function () {
$("div[name='brandsCnt']").hide();
var id=$(this).attr("sid");
$("#"+id).show();
$("div[name='brandsdes']").hide();
var id=$(this).attr("sid");
$('#'+id+'_des').show();
},
function () {

}
);
			
		}
		
	});
}
	
/**
 * 消息滚动功能代码
 */
function msgMove(){
	
	
	// 滚动参数设置
	var ectSlt = ".index_msg";
	var listSlt = ".index_msg_list";
	var elementSlt = ".index_msg_info";
	
	var ect = $(ectSlt);					// 滚动容器
	var listct = $(listSlt, ect);			// 列表容器
	var element = $(elementSlt, listct);	// 元素
	var moveSize = 3;
	var isPause = false;					// 是否暂停状态
	
	if(element.size()<moveSize+1){
		return ;	// 消息小于2条，不进行滚动处理
	}
	
	var length ;							// 滚动总高度
	var total = element.size();				// 总滚动数目
	var index = 0;							// 当前滚动或展示的序号
	var count = 0;							// 移动次数计数
	
	var speedTime = 5000;					// 移动间隔时间
	var speed = 125;						// 滚动时间间隔，影响滚动平滑度
	var movelen = 5;						// 每次移动距离
	var nowTop = 0;
	
	// copy第一条数据至最后
	$(elementSlt+":lt("+moveSize+")", listct).clone().appendTo(listct);
	ect.scrollTop(0);				// 初始化位置
	
	ect.mouseenter(function(){
		isPause = true;
	});
	ect.mouseleave(function(){
		isPause = false;
	});
	
	move();
	
	// 滚动一条消息
	function move(){
		count = 0;
		index++;
		if(index == total+1){	// 当前展示的是最后一个拷贝的，就是第一个了，数据重新初始化
			ect.scrollTop(0);
			index = 1;
		}
		nowTop = ect.scrollTop();
		length = $(elementSlt+":eq("+(index-1)+")", listct).height();	// 滚动总高度
		window.setTimeout(toScroll, speedTime);
	}
	
	// 滚动一小段距离
	function toScroll(){
		if(isPause){	// 暂停中不移动
			window.setTimeout(toScroll, speed);
			return ;
		}
		
		count++;
		
		// 判断滚动完毕
		if(count*movelen >= length){
			ect.scrollTop(nowTop+length);
			move();
		}else{
			ect.scrollTop(nowTop+count*movelen);
			window.setTimeout(toScroll, speed);
		}
	}
}
	

</script>
</head>

<body onload="sideAd()">
<!--页头-->
<!--解决jsp引入乱码 -->
<script src="/pages/2014/js/public.js" type="text/javascript"></script>
<div name="contentType" style="display:none"> <%@ page contentType="text/html; charset=utf-8"%></div>
<div class="head">
  <div class="head-f">
    <div class="width1000 mlr_auto">
      <div class="head-f-gz fr">
      	<a class="block fl ml15 " href="#">关注</a> <span class="icon-head-gz ml5 block fl"></span>
     	 <div class="clear"></div>
     	 <div style="display: none;" class="head-f-ewm">
          	<div class="fl head-f-ewm-img">
            	<img id="top_wx" src=""/>
                <span class="font-color-deep-red block lh22">微信</span>
            </div>
          </div>
          <div class="clear"></div>
        </div>
      <div class="fr head-f-right-nav" style="margin-left:0px;"> 
      	<a href="#" id="bookshelf" target="_blank">我的书架</a> <a class="border-l-solid-gray" id="myBy" target="_blank" href="#">我的消费</a> <a id="myaccount_index" target="_blank" class=" border-l-solid-gray" href="#">我的账户</a> <a id="cooperation" target="_blank" class="border-l-solid-gray" href="#">合作共赢</a> <a id="custService" target="_blank" class="border-l-solid-gray border-r-solid-gray" href="##">意见反馈</a> 
      </div> 
      <div class="icon-small-logo fl"></div>
      <div class="link-bg"> <span class="icon-link-down block fl"></span> <span class="fl">联通旗下网站</span>
        <div class="clear"></div>
        <ul style="display: none;" class="link-bg-ul">
        	<li><a target="_blank" href="http://www.wo.com.cn">沃门户</a></li>
            <li><a target="_blank" href="http://www.10010.com.cn">中国联通网上营业厅</a></li>
        </ul>
      </div>
      <div class="fl ml20">
      	<span id="welcome">你好，请</span>
      	<a  id="userName" class="font-color-deep-red wenzhi" href="#">登录</a>
      </div>
      <a  id="ext" target="_blank" class="block fl ml30 " href="#">免费注册</a>
    </div>
  </div>
  <script>
  		//页头增加二维码
		$(".head-f-gz").mouseover(function(e) {
		   dialogResize("#wxQrcodeIDH", true);
    			$("#wxQrcodeIDH").show();    
        });
		$(".head-f-gz").mouseleave(function(e) {	
           	$("#wxQrcodeIDH").hide();           
        });
		
		//下拉菜单相关网站
		$(".link-bg").mouseover(function(e) {
            $(this).find(".link-bg-ul").show();
        });
		$(".link-bg").mouseleave(function(e) {
            $(this).find(".link-bg-ul").hide();
        });
        
        function submitSearch(){
            var text = $('#conditonArgs').val();
            //var placeholder=$.trim($('#conditonArgs').attr("placeholder"));
            text = $.trim(text);
            if(text == '' || text == undefined || '输入你感兴趣的关键字' == text){
            	return;
            }
            
             //if(text == '' || text == undefined){
           //     showMessage('请输入搜索内容!');
           //if(placeholder==""){
            //  return;
           //}else{
           
           // $('#conditonArgs').val($.trim($('#conditonArgs').attr("placeholder")));
           
          // }
             
           // }else{            
        //      $('#conditonArgs').val($.trim($('#conditonArgs').attr("placeholder")));
           // }

            $('#search2014').submit();
        }
        
        
  </script>

  <div class="head-s">
    <div class="width1000 mlr_auto"> <a href="#"  id="index" class="block fl mt30 head-logo"> </a>
      <div class="fl head-search-bg">
        <form id="search2014" action="/search2014/search_index.action?type=0" method="post" target="_blank">
          <div class="pos_rel fl">
         <input type="text" autocomplete="off" name="conditon" id="conditonArgs" value="输入你感兴趣的关键字" class="head-search-txt">
            <a href="javascript:void(0);"  onclick="deleteContent()" ><span class="icon-fork block"></span></a>
            <div class="mt10 client_select_options wo-pulldown" style="display: none;">
              <ul>
              </ul>
            </div>
            
            <div class="mt10 client_select_options wo-pulldown" style="display: block;">
              <ul id="auto_div">
              </ul>
            </div>
          </div>
        <input type="button" onclick="submitSearch()" value="搜索" class="head-search-btn">
          <div class="clear"></div>
        </form>      
        <div class="head-search-kword" id="searchKword">热门搜索：<a href="#">起风了</a><a href="#">失恋33天</a><a href="#">一吻定情</a></div>
      </div>
      <a href="#" id="client_download" target="_blank">
      <div class="head-client-bg "> <span class="icon-search-client mlr_auto"></span>
        <div class="mt5">客户端</div>
      </div>
      </a>
      <div class="clear"></div>
    </div>
  </div>
  <div class="head-t">
    <div class="width1000 mlr_auto">
      <ul class="head-nav-bg">
        <li><a href="#" id="index_201">首页<span></span></a></li>
        <li><a href="#" id="index_202">出版<span></span></a></li>
        <li><a href="#" id="index_203">原创<span></span></a></li>
        <li><a href="#" id="index_mag">杂志<span></span></a></li>
        <li><a href="#" id="index_aud">听书<span></span></a></li>
        <li><a href="#" id="index_rank">排行榜<span></span></a></li>
        <li><a href="#" id="index_by">包月<span></span></a></li>
        <li><a href="#" id="index_206">免费<span></span></a></li>
        <li><a href="#" id="index_brandzone" >品牌<span></span></a></li>
        <li><a href="#" id="index_actsub">活动<span></span></a></li>
        <div class="clear"></div>
      </ul>
    </div>
  </div>
</div>
<script>
    var type = "201";
    if(type == ""){
    	type = NewGetQueryString("pageindex")
    }
	
	if(type==0){
     var url=window.location.href;
	  if(url.indexOf("index.action")!=-1){
	  $("#index_201").addClass("slt")	  
	  } 
	}else{
    var objname = "#index_" + type;
	//$(".slt").removeClass("slt")
	$(objname).addClass("slt")
	
	}
	
	$(document).ready(function(){
    	var conditon = $("#conditonResult").text();
        $("#conditonArgs").val(conditon);
	});
</script><!--页头结束--> 

<!--腾讯客服-->
<div id="ad_right" style="position: absolute; top: 378px; visibility: visible; left: 1280px;z-index:999;">
	<div class="rt-kf-sidebar">
    	<a href="javascript:void(0);" class="qq-kf" id="qqkfhovera"></a>
        <a href="javascript:void(0);" class="wechat-kf" id="wechatkfhovera"></a>
		<a href="javascript:void(0);" onclick="window.open('http://iread.wo.cn/faq/faqTypeList.action','_blank');">
			<img src="/pages/2014/images/qqwc/yj.png" alt=""/></a>		
        <a href=""><img src="/pages/2014/images/qqwc/back.png" alt=""/></a>	
    </div>
    <div class="kf-hover">
	    <div class="content" id="qqkfhover">
		    <div class="qq-contact-cont">
		    	<dl class="qq-list" style="margin-left:10px;width:90%;">
		        	<dd style="border-bottom:1px dashed #ddd;height:40px;line-height:40px;">
		        	<img src="/pages/2014/images/qqwc/qq.png" alt="" style="position:relative;top:5px;"/><span style="margin-left:12px;">体验师群：</span></dd>
		            <dd>QQ:<a target="_blank" href="http://shang.qq.com/wpa/qunwpa?idkey=8e0ed0cb7c4efac386a08926f7e3f092c0fd9346a1628ffc7685a8209b987886" style="color:#c20200 !important;">42553415</a></dd>
		            <dd>QQ:<a target="_blank" href="http://shang.qq.com/wpa/qunwpa?idkey=a047fffc51903dd550d4a03a3fbd59303265cfe9aef66c0909c70fe6feab1eac" style="color:#c20200 !important;">280198662</a></dd>
		            <dd>QQ:<a target="_blank" href="http://shang.qq.com/wpa/qunwpa?idkey=88675691428e46af7f15b09854324f950fc13631bae85519316a1058f833f6a1" style="color:#c20200 !important;">317117415</a></dd>
		        </dl> 
		    </div>
		</div>
        <div class="wechat-kf-hover" id="wechatkfhover">
        	<img src="/pages/2014/images/qqwc/wechathover.png" alt=""/>
        </div>
    </div>
    <div class="download-phone" id="wxdownloadh">
        <img src="/pages/2014/images/qqwc/khd.png" alt=""/>
    </div>	       
</div> 

<!--腾讯客服 内容-->

<!--主体内容-->
<div class="main pt10">
  <div class="width1000 mlr_auto">
    <div class="main-index-f"> 
      <!--左边分类-->
<div class="main-fl-bg">
        <div class="main-fl-f">         
          <div class="index-l-nav w200">
          	<div class="main-fl-title">
            	<span class="ml20">图书分类</span>		<span class="main-fl-arrows in"></span> 
                <div class="icon-main-title-l"></div>
                <div class="icon-main-title-none"></div>
            </div>
            <ul class="border-l-solid-gray-a main-fl-ul fl f border-r-solid-gray-a b-bottom ml8">
         			            <li class=" slt ">
            	  <a target="_blank" href="/sort2014/sort_index.action?cntType=8&catalogIndex=10003">
              <span class="icon-fl-point"></span>
              <span>都市职业</span></a>         
            	  <a target="_blank" href="/sort2014/sort_index.action?cntType=8&catalogIndex=10001">
              <span class="icon-fl-point"></span>
              <span>玄幻奇幻</span></a>         
	            <div class="clear"></div>
	             </li>
			            <li class="">
            	  <a target="_blank" href="/sort2014/sort_index.action?cntType=8&catalogIndex=10005">
              <span class="icon-fl-point"></span>
              <span>官场职场</span></a>         
            	  <a target="_blank" href="/sort2014/sort_index.action?cntType=8&catalogIndex=10002">
              <span class="icon-fl-point"></span>
              <span>仙侠武侠</span></a>         
	            <div class="clear"></div>
	             </li>
			            <li class=" slt ">
            	  <a target="_blank" href="/sort2014/sort_index.action?cntType=8&catalogIndex=10006">
              <span class="icon-fl-point"></span>
              <span>历史军事</span></a>         
            	  <a target="_blank" href="/sort2014/sort_index.action?cntType=8&catalogIndex=10007">
              <span class="icon-fl-point"></span>
              <span>悬疑灵异</span></a>         
	            <div class="clear"></div>
	             </li>
			            <li class="">
            	  <a target="_blank" href="/sort2014/sort_index.action?cntType=8&catalogIndex=10004">
              <span class="icon-fl-point"></span>
              <span>科幻同人</span></a>         
            	  <a target="_blank" href="/sort2014/sort_index.action?cntType=8&catalogIndex=10008">
              <span class="icon-fl-point"></span>
              <span>网游竞技</span></a>         
	            <div class="clear"></div>
	             </li>
			            <li class=" slt ">
            	  <a target="_blank" href="/sort2014/sort_index.action?cntType=8&catalogIndex=10013">
              <span class="icon-fl-point"></span>
              <span>现代言情</span></a>         
            	  <a target="_blank" href="/sort2014/sort_index.action?cntType=8&catalogIndex=10011">
              <span class="icon-fl-point"></span>
              <span>穿越重生</span></a>         
	            <div class="clear"></div>
	             </li>
			            <li class="">
            	  <a target="_blank" href="/sort2014/sort_index.action?cntType=8&catalogIndex=10010">
              <span class="icon-fl-point"></span>
              <span>古代言情</span></a>         
            	  <a target="_blank" href="/sort2014/sort_index.action?cntType=8&catalogIndex=10009">
              <span class="icon-fl-point"></span>
              <span>幻想言情</span></a>         
	            <div class="clear"></div>
	             </li>
			            <li class=" slt ">
            	  <a target="_blank" href="/sort2014/sort_index.action?cntType=8&catalogIndex=10012">
              <span class="icon-fl-point"></span>
              <span>青春校园</span></a>         
            	  <a target="_blank" href="/sort2014/sort_index.action?cntType=8&catalogIndex=10014">
              <span class="icon-fl-point"></span>
              <span>精品著作</span></a>         
	            <div class="clear"></div>
	             </li>
			            <li class="">
            	  <a target="_blank" href="/sort2014/sort_index.action?cntType=7&catalogIndex=10018">
              <span class="icon-fl-point"></span>
              <span>经管励志</span></a>         
            	  <a target="_blank" href="/sort2014/sort_index.action?cntType=7&catalogIndex=10015">
              <span class="icon-fl-point"></span>
              <span>都市情感</span></a>         
	            <div class="clear"></div>
	             </li>
			            <li class=" slt ">
            	  <a target="_blank" href="/sort2014/sort_index.action?cntType=7&catalogIndex=10016">
              <span class="icon-fl-point"></span>
              <span>传记纪实</span></a>         
            	  <a target="_blank" href="/sort2014/sort_index.action?cntType=7&catalogIndex=10020">
              <span class="icon-fl-point"></span>
              <span>人文社科</span></a>         
	            <div class="clear"></div>
		            </li>
              
              <div class="clear"></div>
            </ul>
            <ul style="display: none;" class="main-fl-ul fl f   border-r-solid-gray-a b-bottom">
            <li class=" slt ">
                <a target="_blank" href="/sort2014/sort_index.action?cntType=7&catalogIndex=10021">
              <span class="icon-fl-point"></span>
              <span>时尚生活</span></a>         
                <a target="_blank" href="/sort2014/sort_index.action?cntType=7&catalogIndex=10019">
              <span class="icon-fl-point"></span>
              <span>影视悬疑</span></a>         
            <div class="clear"></div>
             </li>          
            <li class="">
                <a target="_blank" href="/sort2014/sort_index.action?cntType=7&catalogIndex=10017">
              <span class="icon-fl-point"></span>
              <span>职场百态</span></a>         
                <a target="_blank" href="/sort2014/sort_index.action?cntType=8&catalogIndex=10022">
              <span class="icon-fl-point"></span>
              <span>古典名著</span></a>         
            <div class="clear"></div>
             </li>          
            <li class=" slt ">
                <a target="_blank" href="/sort2014/sort_index.action?cntType=7&catalogIndex=10024">
              <span class="icon-fl-point"></span>
              <span>成长教育</span></a>         
                <a target="_blank" href="/sort2014/sort_index.action?cntType=7&catalogIndex=10023">
              <span class="icon-fl-point"></span>
              <span>国学经典</span></a>         
            <div class="clear"></div>
             </li>          
            <li class="">
                <a target="_blank" href="/sort2014/sort_index.action?cntType=7&catalogIndex=10025">
              <span class="icon-fl-point"></span>
              <span>其他更多</span></a>         
            <div class="clear"></div>
          </li>         
              <div class="clear"></div>
            </ul>
            
            
       <ul style="display: none;" class="main-fl-ul fl f border-r-solid-gray-a b-bottom">
              <div class="clear"></div>
            </ul>
            
            
            <div class="clear"></div>
          </div>
        </div>
        <div class="main-fl-s">
          
          <div class="index-l-nav w200">
          	<div class="main-fl-title">
            	<span class="ml20">杂志分类</span>		<span class="main-fl-arrows in"></span> 
                <div class="icon-main-title-l"></div>
                <div class="icon-main-title-none"></div>
            </div>
            <ul class="border-l-solid-gray-a main-fl-ul  fl ml8 f2  border-r-solid-gray-a b-bottom">
            <li class=" slt ">
                  <a target="_blank" href="/sort2014/sort_index.action?cntType=3&catalogIndex=119">
              <span class="icon-fl-point"></span>
              <span>文学</span></a>         
                
             
                  <a target="_blank" href="/sort2014/sort_index.action?cntType=3&catalogIndex=113">
              <span class="icon-fl-point"></span>
              <span>故事</span></a>         
            <div class="clear"></div>
             </li>          
                
             
            <li class="">
                  <a target="_blank" href="/sort2014/sort_index.action?cntType=3&catalogIndex=107">
              <span class="icon-fl-point"></span>
              <span>生活</span></a>         
                
             
                  <a target="_blank" href="/sort2014/sort_index.action?cntType=3&catalogIndex=101">
              <span class="icon-fl-point"></span>
              <span>时尚</span></a>         
            <div class="clear"></div>
             </li>          
                
             
            <li class=" slt ">
                  <a target="_blank" href="/sort2014/sort_index.action?cntType=3&catalogIndex=104">
              <span class="icon-fl-point"></span>
              <span>娱乐</span></a>         
                
             
                  <a target="_blank" href="/sort2014/sort_index.action?cntType=3&catalogIndex=102">
              <span class="icon-fl-point"></span>
              <span>时事</span></a>         
            <div class="clear"></div>
             </li>          
                
             
              <div class="clear"></div>
            </ul>
            <ul style="display: none;" class="main-fl-ul fl f2  border-r-solid-gray-a b-bottom">
          
               
            <li class=" slt ">
                <a target="_blank" href="/sort2014/sort_index.action?cntType=3&catalogIndex=103">
              <span class="icon-fl-point"></span>
              <span>财经</span></a>         
                
          
               
                <a target="_blank" href="/sort2014/sort_index.action?cntType=3&catalogIndex=121">
              <span class="icon-fl-point"></span>
              <span>管理</span></a>         
            <div class="clear"></div>
             </li>          
                
          
               
            <li class="">
                <a target="_blank" href="/sort2014/sort_index.action?cntType=3&catalogIndex=112">
              <span class="icon-fl-point"></span>
              <span>军事</span></a>         
                
          
               
                <a target="_blank" href="/sort2014/sort_index.action?cntType=3&catalogIndex=108">
              <span class="icon-fl-point"></span>
              <span>健康</span></a>         
            <div class="clear"></div>
             </li>          
                
          
               
            <li class=" slt ">
                <a target="_blank" href="/sort2014/sort_index.action?cntType=3&catalogIndex=115">
              <span class="icon-fl-point"></span>
              <span>数码</span></a>         
                
          
               
                <a target="_blank" href="/sort2014/sort_index.action?cntType=3&catalogIndex=111">
              <span class="icon-fl-point"></span>
              <span>旅游</span></a>         
            <div class="clear"></div>
             </li>          
                
             
             
             
              <div class="clear"></div>
            </ul>
                 <ul style="display: none;" class="main-fl-ul fl f2 border-r-solid-gray-a b-bottom">
          
               
            <li class=" slt ">
                <a target="_blank" href="/sort2014/sort_index.action?cntType=3&catalogIndex=106">
              <span class="icon-fl-point"></span>
              <span>汽车</span></a>         
                
          
               
                <a target="_blank" href="/sort2014/sort_index.action?cntType=3&catalogIndex=110">
              <span class="icon-fl-point"></span>
              <span>教育</span></a>         
            <div class="clear"></div>
             </li>          
                
          
               
            <li class="">
                <a target="_blank" href="/sort2014/sort_index.action?cntType=3&catalogIndex=117">
              <span class="icon-fl-point"></span>
              <span>社科</span></a>         
                
          
               
                <a target="_blank" href="/sort2014/sort_index.action?cntType=3&catalogIndex=105">
              <span class="icon-fl-point"></span>
              <span>体育</span></a>         
            <div class="clear"></div>
             </li>          
                
          
               
            <li class=" slt ">
                <a target="_blank" href="/sort2014/sort_index.action?cntType=3&catalogIndex=114">
              <span class="icon-fl-point"></span>
              <span>科普</span></a>         
                
          
               
                <a target="_blank" href="/sort2014/sort_index.action?cntType=3&catalogIndex=116">
              <span class="icon-fl-point"></span>
              <span>人文</span></a>         
            <div class="clear"></div>
             </li>          
                
             
             
             
              <div class="clear"></div>
            </ul>
            
                             <ul style="display: none;" class="main-fl-ul fl  f2 border-r-solid-gray-a b-bottom">
          
               
            <li class=" slt ">
                <a target="_blank" href="/sort2014/sort_index.action?cntType=3&catalogIndex=120">
              <span class="icon-fl-point"></span>
              <span>行业</span></a>         
                
          
               
                <a target="_blank" href="/sort2014/sort_index.action?cntType=3&catalogIndex=109">
              <span class="icon-fl-point"></span>
              <span>奢华</span></a>         
            <div class="clear"></div>
             </li>          
                
          
               
            <li class="">
                <a target="_blank" href="/sort2014/sort_index.action?cntType=3&catalogIndex=118">
              <span class="icon-fl-point"></span>
              <span>其他</span></a>         
            <div class="clear"></div>
          </li>         
                
             
             
             
              <div class="clear"></div>
            </ul>
            
            <div class="clear"></div>
          </div>
        </div>
        <div class="main-fl-t">
          
          <div class="index-l-nav w200">
          <div class="main-fl-title">
            	<span class="ml20">听书分类</span>		<span class="main-fl-arrows in"></span> 
                <div class="icon-main-title-l"></div>
                <div class="icon-main-title-none"></div>
            </div>
            <ul class="border-l-solid-gray-a main-fl-ul fl  f2 ml8 border-r-solid-gray-a b-bottom">
            <li class=" slt ">
                <a target="_blank" href="/sort2014/sort_index.action?cntType=5&catalogIndex=1060">
              <span class="icon-fl-point"></span>
              <span>玄幻武侠</span></a>         
                <a target="_blank" href="/sort2014/sort_index.action?cntType=5&catalogIndex=1061">
              <span class="icon-fl-point"></span>
              <span>都市言情</span></a>         
            <div class="clear"></div>
             </li>          
            <li class="">
                <a target="_blank" href="/sort2014/sort_index.action?cntType=5&catalogIndex=1062">
              <span class="icon-fl-point"></span>
              <span>历史军事</span></a>         
                <a target="_blank" href="/sort2014/sort_index.action?cntType=5&catalogIndex=1063">
              <span class="icon-fl-point"></span>
              <span>恐怖悬疑</span></a>         
            <div class="clear"></div>
             </li>          
            <li class=" slt ">
                <a target="_blank" href="/sort2014/sort_index.action?cntType=5&catalogIndex=1064">
              <span class="icon-fl-point"></span>
              <span>经典文学</span></a>         
                <a target="_blank" href="/sort2014/sort_index.action?cntType=5&catalogIndex=1065">
              <span class="icon-fl-point"></span>
              <span>相声评书</span></a>         
            <div class="clear"></div>
             </li>          
              <div class="clear"></div>
            </ul>
            <ul style="display: none;" class="main-fl-ul fl  f2 border-r-solid-gray-a b-bottom">
            <li class=" slt ">
                <a target="_blank" href="/sort2014/sort_index.action?cntType=5&catalogIndex=1066">
              <span class="icon-fl-point"></span>
              <span>幽默短篇</span></a>         
                <a target="_blank" href="/sort2014/sort_index.action?cntType=5&catalogIndex=1067">
              <span class="icon-fl-point"></span>
              <span>综艺娱乐</span></a>         
            <div class="clear"></div>
             </li>          
            <li class="">
                <a target="_blank" href="/sort2014/sort_index.action?cntType=5&catalogIndex=1068">
              <span class="icon-fl-point"></span>
              <span>少儿读物</span></a>         
                <a target="_blank" href="/sort2014/sort_index.action?cntType=5&catalogIndex=1069">
              <span class="icon-fl-point"></span>
              <span>经管励志</span></a>         
            <div class="clear"></div>
             </li>          
            <li class=" slt ">
                <a target="_blank" href="/sort2014/sort_index.action?cntType=5&catalogIndex=1059">
              <span class="icon-fl-point"></span>
              <span>教育科普</span></a>         
            <div class="clear"></div>
          </li>         
              <div class="clear"></div>
            </ul>   
  
            <div class="clear"></div>
          </div>
        </div>
      </div>      <!--左边分类end--> 
      
      <!--中间-->
      <div class="main-c-bg" style="position: absolute; left: 190px;">
      
       <!--大banner start-->
        <div class="index_banner">
          <div class="banner_list">
            <ul>
							       		<li id="img_bar_li_3" class='slt'  ><a href="http://iread.wo.com.cn/actsub/selectActivityInfo.action?activityID=3576" target="_blank" ><img src="http://iread.wo.com.cn/specialarea/20170511101109_banner.jpg" alt=""></a></li>           
							       		<li id="img_bar_li_4" class=''  ><a href="http://iread.wo.com.cn/actsub/selectActivityInfo.action?activityID=2457" target="_blank" ><img src="http://iread.wo.com.cn/specialarea/20160526163420_banner.jpg" alt=""></a></li>           
            </ul>
          </div>
          <div class="banner_bar">
            <ul>
        							<li id="bar_li_3"  >
                          
        							<li id="bar_li_4"  >
                          
            </ul>
          </div>
        </div>
        <div> 
	        <a href="/wofun/wofundetail.action" target="_blank" class="index-s-banner ">
	        <img src="http://iread.wo.com.cn/specialarea/20141031052005.jpg" alt=""></a>	        




          <div class="clear"></div>
        </div>        
         <!--大banner end-->
        <div class="main-c-t-bg">
          <div class="main-c-tab-bg mt10">
          
          <!--固定3个栏目  只拿3个-->
            <div  class="main-c-tab-cell   slt" name="showCatindex" hid="0_cat" onclick="showBook('0_cat',this)">新书</div>
            <div  class="main-c-tab-cell   " name="showCatindex" hid="1_cat" onclick="showBook('1_cat',this)">重磅全免</div>
            <div  class="main-c-tab-cell   " name="showCatindex" hid="2_cat" onclick="showBook('2_cat',this)">主编力荐</div>
            <div class="clear"></div>
          </div>
          <!--固定3个栏目 书 只拿2本书-->
          <div class="main-c-t-book" name="catindex" id="0_cat" >                  
            <div class="main-c-t-cell-bg border-r-solid-gray" > <a href="/contentdetail/detail.action?cntindex=476177&catid=106634"  target="_blank" class="book-img ml10"> <img src="http://iread.wo.com.cn/cnt/image/476/10476177/cover_176.jpg" alt=""> </a>
              <div class="fl ml10 book-right"> <a href="/contentdetail/detail.action?cntindex=476177&catid=106634" target="_blank" class="book-title">重生农媳逆袭</a>
                <div class="mt10"> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                 <span class="star_none"></span>
                  <div class="clear"></div>
                </div>
                <div class="book-txt mt10"> <span class="font-color-gray">作者：</span><span > 懒百合</span></div>
                <div class="book-txt mt10"> <span class="font-color-gray">简介：</span>连翘一睁眼，发现自己重回到了1987年，...</div>
              </div>
              <div class="clear"></div>
            </div>
            <div class="main-c-t-cell-bg " > <a href="/contentdetail/detail.action?cntindex=396521&catid=106634"  target="_blank" class="book-img ml10"> <img src="http://iread.wo.com.cn/cnt/image/396/10396521/cover_176.jpg" alt=""> </a>
              <div class="fl ml10 book-right"> <a href="/contentdetail/detail.action?cntindex=396521&catid=106634" target="_blank" class="book-title">重生之绝世武神</a>
                <div class="mt10"> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                 <span class="star_none"></span>
                  <div class="clear"></div>
                </div>
                <div class="book-txt mt10"> <span class="font-color-gray">作者：</span><span > 风一刀</span></div>
                <div class="book-txt mt10"> <span class="font-color-gray">简介：</span>绝代武神杨腾因得到一件帝器，遭人陷害引来...</div>
              </div>
              <div class="clear"></div>
            </div>
            <div class="clear"></div>
          </div>
          <!--固定3个栏目 书 只拿2本书-->
          <div class="main-c-t-book" name="catindex" id="1_cat" style="display:none">                  
            <div class="main-c-t-cell-bg border-r-solid-gray" > <a href="/contentdetail/detail.action?cntindex=480388&catid=107915"  target="_blank" class="book-img ml10"> <img src="http://iread.wo.com.cn/cnt/image/480/10480388/cover_176.jpg" alt=""> </a>
              <div class="fl ml10 book-right"> <a href="/contentdetail/detail.action?cntindex=480388&catid=107915" target="_blank" class="book-title">农夫传奇</a>
                <div class="mt10"> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                 <span class="star_none"></span>
                  <div class="clear"></div>
                </div>
                <div class="book-txt mt10"> <span class="font-color-gray">作者：</span><span > 关汉时</span></div>
                <div class="book-txt mt10"> <span class="font-color-gray">简介：</span>毕业后找不到工作？不要紧，回家种田吧！他...</div>
              </div>
              <div class="clear"></div>
            </div>
            <div class="main-c-t-cell-bg " > <a href="/contentdetail/detail.action?cntindex=476805&catid=107915"  target="_blank" class="book-img ml10"> <img src="http://iread.wo.com.cn/cnt/image/476/10476805/cover_176.jpg" alt=""> </a>
              <div class="fl ml10 book-right"> <a href="/contentdetail/detail.action?cntindex=476805&catid=107915" target="_blank" class="book-title">都市最强装逼系统</a>
                <div class="mt10"> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                 <span class="star_none"></span>
                  <div class="clear"></div>
                </div>
                <div class="book-txt mt10"> <span class="font-color-gray">作者：</span><span > 必火</span></div>
                <div class="book-txt mt10"> <span class="font-color-gray">简介：</span>叶秋重生获得装逼系统，当最萌萝莉的奶爸，...</div>
              </div>
              <div class="clear"></div>
            </div>
            <div class="clear"></div>
          </div>
          <!--固定3个栏目 书 只拿2本书-->
          <div class="main-c-t-book" name="catindex" id="2_cat" style="display:none">                  
            <div class="main-c-t-cell-bg border-r-solid-gray" > <a href="/contentdetail/detail.action?cntindex=454082&catid=106626"  target="_blank" class="book-img ml10"> <img src="http://iread.wo.com.cn/cnt/image/454/10454082/cover_176.jpg" alt=""> </a>
              <div class="fl ml10 book-right"> <a href="/contentdetail/detail.action?cntindex=454082&catid=106626" target="_blank" class="book-title">最强小民工</a>
                <div class="mt10"> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                 <span class="star_none"></span>
                  <div class="clear"></div>
                </div>
                <div class="book-txt mt10"> <span class="font-color-gray">作者：</span><span > 夜深自呓</span></div>
                <div class="book-txt mt10"> <span class="font-color-gray">简介：</span>代号龙王的传奇男人，自监狱中走出当起了小...</div>
              </div>
              <div class="clear"></div>
            </div>
            <div class="main-c-t-cell-bg " > <a href="/contentdetail/detail.action?cntindex=449822&catid=106626"  target="_blank" class="book-img ml10"> <img src="http://iread.wo.com.cn/cnt/image/449/10449822/cover_176.jpg" alt=""> </a>
              <div class="fl ml10 book-right"> <a href="/contentdetail/detail.action?cntindex=449822&catid=106626" target="_blank" class="book-title">首席情深：豪门第一夫人</a>
                <div class="mt10"> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                 <span class="star_none"></span>
                  <div class="clear"></div>
                </div>
                <div class="book-txt mt10"> <span class="font-color-gray">作者：</span><span > 夏寒烟</span></div>
                <div class="book-txt mt10"> <span class="font-color-gray">简介：</span>叶氏破产，朝夕之间，她一无所有。权势滔天...</div>
              </div>
              <div class="clear"></div>
            </div>
            <div class="clear"></div>
          </div>
          
          
        </div>
      </div>
      
      <!--中间end--> 
      
      <!--右边-->
      
      <div class="main-r-bg">
        <div class="main-r-border">
          <div class="main-r-f-bg pb10" id="readAndNologin"  style="height:206px;">
         
            
            
            </div>
         <!-- 站内公告 start-->
          <div class="main-r-msg-bg mt10">
            <div class="main-r-msg-title">站内公告<a class="font-size-12 font-color-deep-red fr" href="/notice/list.action" target="_blank">更多&gt;&gt;</a>
              <div class="clear"></div>
            </div>
            <div class="index_msg">
            	<div class="index_msg_list">
       			</div>  
       		</div>  
          </div>
        </div>
        <!-- 站内公告 end-->
        
        <div class="main-c-t-bg">
          <div class="main-c-tab-bg mt10">
            <div class="main-c-tab-cell slt" hid="wx" onclick="show('wx',this)" name="wx"> 微信关注 </div>
            <div class="clear"></div>
          </div>
         <div class="main-c-t-down" id="wx" name="down">
            <div class="main-msg-ewm1 mlr_auto">  <img  id="clienth" class="qrcode-border" alt="" src="/pages/2014/images/ercode/wowx.jpg"></div>
            <div class="ac font-size-16 font-color-black mt15">微信关注 </div>
            <div class="ac font-size-12 font-color-gray mt10">扫一扫，体验更多精选内容分享。</div>
          </div>
        </div>
      </div>
      <!--右边end--> 
     <div class="clear"></div>
    </div>
    
    <!--分类板块-->

       <div>
      <div class="mod-a-title mt10 110,468">
        <div class="mod-a-title-txt fl"> 
        <a name="cataliasTag" target="_blank" href="/index.action?pageindex=202">精品专区</a>
          <div class="clear"></div>
        </div>
        	<a name="cataliasTag" target="_blank" class="font-color-deep-red font-color-deep-red font-size-12 fr mt5" href="/index.action?pageindex=202">          
       			进入
精品专区
       			频道&gt;&gt;
        	</a>
        <div class="clear"></div>
      </div>
      <div class="mt10">
        <div class="mod-a-l-bg" id="index-nspd_110468">
        
          <div class="mod-a-banner">
            <div class="mod-a-banner-list">
              <ul>
			               
								              <li class="img_bar_li_0 
								              	slt" id="">
								              		<a href="/index.action?pageindex=202" target="_blank"> <img src="http://iread.wo.com.cn/specialarea/20141021102244.jpg" alt=""></a>
								              </li>              
              </ul>
            </div>
            <div class="mod-a-banner-bar">
              <ul>
				                			<li class="" id="bar_li_0" style="display:none;"></li>
                
              </ul>
            </div>
          </div>
          <div class=" mlr10 mt5 pb5 border-b-solid-gray"></div>
	          <div class=""> 
				          		 <ul style="font-size:14px;line-height:30px;">
				          		 	<li>
								      <a class="bklist_bk" target="_blank" href="/contentdetail/detail.action?cntindex=485399&catid=110468">
								      	<span class="ran_index_flag">1</span>
								      	<div style="font-family:微软雅黑,宋体,Arial,Helvetica,sans-serif;white-space:nowrap;word-break:break-all;word-wrap:break-word;">人民的名义</div>
								      </a>
							      	</li>
							     </ul>
				          		 <ul style="font-size:14px;line-height:30px;">
				          		 	<li>
								      <a class="bklist_bk" target="_blank" href="/contentdetail/detail.action?cntindex=480206&catid=110468">
								      	<span class="ran_index_flag">2</span>
								      	<div style="font-family:微软雅黑,宋体,Arial,Helvetica,sans-serif;white-space:nowrap;word-break:break-all;word-wrap:break-word;">三生三世十里桃花</div>
								      </a>
							      	</li>
							     </ul>
				          		 <ul style="font-size:14px;line-height:30px;">
				          		 	<li>
								      <a class="bklist_bk" target="_blank" href="/contentdetail/detail.action?cntindex=10017014&catid=110468">
								      	<span class="ran_index_flag">3</span>
								      	<div style="font-family:微软雅黑,宋体,Arial,Helvetica,sans-serif;white-space:nowrap;word-break:break-all;word-wrap:break-word;">鬼吹灯</div>
								      </a>
							      	</li>
							     </ul>
				          		 <ul style="font-size:14px;line-height:30px;">
				          		 	<li>
								      <a class="bklist_bk" target="_blank" href="/contentdetail/detail.action?cntindex=392043&catid=110468">
								      	<span class="ran_index_flag">4</span>
								      	<div style="font-family:微软雅黑,宋体,Arial,Helvetica,sans-serif;white-space:nowrap;word-break:break-all;word-wrap:break-word;">上瘾五百年</div>
								      </a>
							      	</li>
							     </ul>
				          		 <ul style="font-size:14px;line-height:30px;">
				          		 	<li>
								      <a class="bklist_bk" target="_blank" href="/contentdetail/detail.action?cntindex=324110&catid=110468">
								      	<span class="ran_index_flag">5</span>
								      	<div style="font-family:微软雅黑,宋体,Arial,Helvetica,sans-serif;white-space:nowrap;word-break:break-all;word-wrap:break-word;">U．I．S校园日志</div>
								      </a>
							      	</li>
							     </ul>
	             </div>
        	</div>
        	<div class="fl">
				         		<div class="ml20 pb10">
							   <div class="mod-a-cell-bg mlr20" name="0" ss="10"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=480206&catid=110468"><img src="http://iread.wo.com.cn/cnt/image/480/10480206/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="三生三世十里桃花" href="/contentdetail/detail.action?cntindex=480206&catid=110468">三生三世十里桃花</a> <span class="block"> 唐七</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="1" ss="10"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=480189&catid=110468"><img src="http://iread.wo.com.cn/cnt/image/480/10480189/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="为什么精英都是清单控" href="/contentdetail/detail.action?cntindex=480189&catid=110468">为什么精英都是清单控</a> <span class="block"> 【美】宝拉·里佐/郑焕升</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="2" ss="10"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=479733&catid=110468"><img src="http://iread.wo.com.cn/cnt/image/479/10479733/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="辽东轶闻手记" href="/contentdetail/detail.action?cntindex=479733&catid=110468">辽东轶闻手记</a> <span class="block"> 叶遁</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="3" ss="10"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=479661&catid=110468"><img src="http://iread.wo.com.cn/cnt/image/479/10479661/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="每个人都是一个宇宙" href="/contentdetail/detail.action?cntindex=479661&catid=110468">每个人都是一个宇宙</a> <span class="block"> 周国平</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="4" ss="10"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=479654&catid=110468"><img src="http://iread.wo.com.cn/cnt/image/479/10479654/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="军锋" href="/contentdetail/detail.action?cntindex=479654&catid=110468">军锋</a> <span class="block"> 冷海</span>
							   </div>
            					<div class="clear"></div>
          						</div>
				         		<div class="ml20 pb10">
							   <div class="mod-a-cell-bg mlr20" name="5" ss="10"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=479613&catid=110468"><img src="http://iread.wo.com.cn/cnt/image/479/10479613/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="那些年我们错过的勇气" href="/contentdetail/detail.action?cntindex=479613&catid=110468">那些年我们错过的勇气</a> <span class="block"> 林蓠</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="6" ss="10"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=479609&catid=110468"><img src="http://iread.wo.com.cn/cnt/image/479/10479609/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="后来，你都如何回忆我" href="/contentdetail/detail.action?cntindex=479609&catid=110468">后来，你都如何回忆我</a> <span class="block"> 那时迷离</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="7" ss="10"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=475575&catid=110468"><img src="http://iread.wo.com.cn/cnt/image/475/10475575/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="生物演化与人类未来" href="/contentdetail/detail.action?cntindex=475575&catid=110468">生物演化与人类未来</a> <span class="block"> 殷鸿福</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="8" ss="10"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=475572&catid=110468"><img src="http://iread.wo.com.cn/cnt/image/475/10475572/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="一切有情，都无挂碍" href="/contentdetail/detail.action?cntindex=475572&catid=110468">一切有情，都无挂碍</a> <span class="block"> 马文戈</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="9" ss="10"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=474207&catid=110468"><img src="http://iread.wo.com.cn/cnt/image/474/10474207/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="因为你，我爱上这世界" href="/contentdetail/detail.action?cntindex=474207&catid=110468">因为你，我爱上这世界</a> <span class="block"> 蜃公子</span>
							   </div>
            					<div class="clear"></div>
          						</div>
             </div>
        	<div class="clear"></div>
      	</div>
    		<div class="clear"></div>
    	</div>

       <div>
      <div class="mod-a-title mt10 106,635">
        <div class="mod-a-title-txt fl"> 
        <span class="icon-a-mod nan mt10 mr10">
        
        </span>
        <a name="cataliasTag" target="_blank" href="index.action?pageindex=204#">男生频道</a>
          <div class="clear"></div>
        </div>
        	<a name="cataliasTag" target="_blank" class="font-color-deep-red font-color-deep-red font-size-12 fr mt5" href="index.action?pageindex=204#">          
       			进入
       				
       				男生
       				频道&gt;&gt;
        	</a>
        <div class="clear"></div>
      </div>
      <div class="mt10">
        <div class="mod-a-l-bg" id="index-nspd_106635">
        
          <div class="mod-a-banner">
            <div class="mod-a-banner-list">
              <ul>
			               
								              <li class="img_bar_li_0 
								              	slt" id="">
								              		<a href="http://iread.wo.com.cn/stacks/getContentDetail.action?cntindex=315128&discountindex=" target="_blank"> <img src="http://iread.wo.com.cn/specialarea/20140915091635.jpg" alt=""></a>
								              </li>              
								              <li class="img_bar_li_1 
								              	" id="">
								              		<a href="http://iread.wo.com.cn/stacks/getContentDetail.action?cntindex=297352&discountindex=" target="_blank"> <img src="http://iread.wo.com.cn/specialarea/20140915091720.jpg" alt=""></a>
								              </li>              
              </ul>
            </div>
            <div class="mod-a-banner-bar">
              <ul>
				                			<li class="" id="bar_li_0" style="display:none;"></li>
				                			<li class="" id="bar_li_1" style="display:none;"></li>
                
              </ul>
            </div>
          </div>
          <div class=" mlr10 mt5 pb5 border-b-solid-gray"></div>
	          <div class=""> 
				          		 <ul style="font-size:14px;line-height:30px;">
				          		 	<li>
								      <a class="bklist_bk" target="_blank" href="/contentdetail/detail.action?cntindex=432198&catid=106635">
								      	<span class="ran_index_flag">1</span>
								      	<div style="font-family:微软雅黑,宋体,Arial,Helvetica,sans-serif;white-space:nowrap;word-break:break-all;word-wrap:break-word;">择天记</div>
								      </a>
							      	</li>
							     </ul>
				          		 <ul style="font-size:14px;line-height:30px;">
				          		 	<li>
								      <a class="bklist_bk" target="_blank" href="/contentdetail/detail.action?cntindex=218333&catid=106635">
								      	<span class="ran_index_flag">2</span>
								      	<div style="font-family:微软雅黑,宋体,Arial,Helvetica,sans-serif;white-space:nowrap;word-break:break-all;word-wrap:break-word;">校花的贴身高手</div>
								      </a>
							      	</li>
							     </ul>
				          		 <ul style="font-size:14px;line-height:30px;">
				          		 	<li>
								      <a class="bklist_bk" target="_blank" href="/contentdetail/detail.action?cntindex=474517&catid=106635">
								      	<span class="ran_index_flag">3</span>
								      	<div style="font-family:微软雅黑,宋体,Arial,Helvetica,sans-serif;white-space:nowrap;word-break:break-all;word-wrap:break-word;">无敌悍民</div>
								      </a>
							      	</li>
							     </ul>
				          		 <ul style="font-size:14px;line-height:30px;">
				          		 	<li>
								      <a class="bklist_bk" target="_blank" href="/contentdetail/detail.action?cntindex=405698&catid=106635">
								      	<span class="ran_index_flag">4</span>
								      	<div style="font-family:微软雅黑,宋体,Arial,Helvetica,sans-serif;white-space:nowrap;word-break:break-all;word-wrap:break-word;">最强兵王</div>
								      </a>
							      	</li>
							     </ul>
				          		 <ul style="font-size:14px;line-height:30px;">
				          		 	<li>
								      <a class="bklist_bk" target="_blank" href="/contentdetail/detail.action?cntindex=474506&catid=106635">
								      	<span class="ran_index_flag">5</span>
								      	<div style="font-family:微软雅黑,宋体,Arial,Helvetica,sans-serif;white-space:nowrap;word-break:break-all;word-wrap:break-word;">绝品护花兵王</div>
								      </a>
							      	</li>
							     </ul>
	             </div>
        	</div>
        	<div class="fl">
				         		<div class="ml20 pb10">
							   <div class="mod-a-cell-bg mlr20" name="0" ss="4"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=292928&catid=106635"><img src="http://iread.wo.com.cn/cnt/image/292/10292928/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="绝品天医" href="/contentdetail/detail.action?cntindex=292928&catid=106635">绝品天医</a> <span class="block"> 叶天南</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="1" ss="4"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=339538&catid=106635"><img src="http://iread.wo.com.cn/cnt/image/339/10339538/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="无敌天下" href="/contentdetail/detail.action?cntindex=339538&catid=106635">无敌天下</a> <span class="block"> 神见</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="2" ss="4"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=293185&catid=106635"><img src="http://iread.wo.com.cn/cnt/image/293/10293185/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="重生之官场风流" href="/contentdetail/detail.action?cntindex=293185&catid=106635">重生之官场风流</a> <span class="block"> 香烟盒子</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="3" ss="4"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=476690&catid=106635"><img src="http://iread.wo.com.cn/cnt/image/476/10476690/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="修仙进行中" href="/contentdetail/detail.action?cntindex=476690&catid=106635">修仙进行中</a> <span class="block"> 暗夜泠风</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="4" ss="4"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=293773&catid=106635"><img src="http://iread.wo.com.cn/cnt/image/293/10293773/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="权柄" href="/contentdetail/detail.action?cntindex=293773&catid=106635">权柄</a> <span class="block"> 三戒大师</span>
							   </div>
            					<div class="clear"></div>
          						</div>
				         		<div class="ml20 pb10">
							   <div class="mod-a-cell-bg mlr20" name="5" ss="4"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=438055&catid=106635"><img src="http://iread.wo.com.cn/cnt/image/438/10438055/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="巅峰小农民" href="/contentdetail/detail.action?cntindex=438055&catid=106635">巅峰小农民</a> <span class="block"> 白菜汤</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="6" ss="4"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=391478&catid=106635"><img src="http://iread.wo.com.cn/cnt/image/391/10391478/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="绝世药皇" href="/contentdetail/detail.action?cntindex=391478&catid=106635">绝世药皇</a> <span class="block"> 永远天涯</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="7" ss="4"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=229916&catid=106635"><img src="http://iread.wo.com.cn/cnt/image/229/10229916/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="官路修行" href="/contentdetail/detail.action?cntindex=229916&catid=106635">官路修行</a> <span class="block"> 蔡晋</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="8" ss="4"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=218339&catid=106635"><img src="http://iread.wo.com.cn/cnt/image/218/10218339/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="遮天" href="/contentdetail/detail.action?cntindex=218339&catid=106635">遮天</a> <span class="block"> 辰东</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="9" ss="4"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=430668&catid=106635"><img src="http://iread.wo.com.cn/cnt/image/430/10430668/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="铁血抗日" href="/contentdetail/detail.action?cntindex=430668&catid=106635">铁血抗日</a> <span class="block"> 残简</span>
							   </div>
            					<div class="clear"></div>
          						</div>
             </div>
        	<div class="clear"></div>
      	</div>
    		<div class="clear"></div>
    	</div>

       <div>
      <div class="mod-a-title mt10 106,781">
        <div class="mod-a-title-txt fl"> 
          <span class="icon-a-mod nv mt10 mr10">
        </span>
        <a name="cataliasTag" target="_blank" href="/index.action?pageindex=205#">女生频道</a>
          <div class="clear"></div>
        </div>
        	<a name="cataliasTag" target="_blank" class="font-color-deep-red font-color-deep-red font-size-12 fr mt5" href="/index.action?pageindex=205#">          
       			进入
       				女生
       				
       				频道&gt;&gt;
        	</a>
        <div class="clear"></div>
      </div>
      <div class="mt10">
        <div class="mod-a-l-bg" id="index-nspd_106781">
        
          <div class="mod-a-banner">
            <div class="mod-a-banner-list">
              <ul>
			               
								              <li class="img_bar_li_0 
								              	slt" id="">
								              		<a href="/index.action?pageindex=205#" target="_blank"> <img src="http://iread.wo.com.cn/specialarea/20140912113758.jpg" alt=""></a>
								              </li>              
              </ul>
            </div>
            <div class="mod-a-banner-bar">
              <ul>
				                			<li class="" id="bar_li_0" style="display:none;"></li>
                
              </ul>
            </div>
          </div>
          <div class=" mlr10 mt5 pb5 border-b-solid-gray"></div>
	          <div class=""> 
				          		 <ul style="font-size:14px;line-height:30px;">
				          		 	<li>
								      <a class="bklist_bk" target="_blank" href="/contentdetail/detail.action?cntindex=449913&catid=106781">
								      	<span class="ran_index_flag">1</span>
								      	<div style="font-family:微软雅黑,宋体,Arial,Helvetica,sans-serif;white-space:nowrap;word-break:break-all;word-wrap:break-word;">闪婚掠爱：帝少宠妻入骨</div>
								      </a>
							      	</li>
							     </ul>
				          		 <ul style="font-size:14px;line-height:30px;">
				          		 	<li>
								      <a class="bklist_bk" target="_blank" href="/contentdetail/detail.action?cntindex=315076&catid=106781">
								      	<span class="ran_index_flag">2</span>
								      	<div style="font-family:微软雅黑,宋体,Arial,Helvetica,sans-serif;white-space:nowrap;word-break:break-all;word-wrap:break-word;">总裁，我要离婚</div>
								      </a>
							      	</li>
							     </ul>
				          		 <ul style="font-size:14px;line-height:30px;">
				          		 	<li>
								      <a class="bklist_bk" target="_blank" href="/contentdetail/detail.action?cntindex=438441&catid=106781">
								      	<span class="ran_index_flag">3</span>
								      	<div style="font-family:微软雅黑,宋体,Arial,Helvetica,sans-serif;white-space:nowrap;word-break:break-all;word-wrap:break-word;">BOSS来袭：娇妻躺下，别闹！</div>
								      </a>
							      	</li>
							     </ul>
				          		 <ul style="font-size:14px;line-height:30px;">
				          		 	<li>
								      <a class="bklist_bk" target="_blank" href="/contentdetail/detail.action?cntindex=438064&catid=106781">
								      	<span class="ran_index_flag">4</span>
								      	<div style="font-family:微软雅黑,宋体,Arial,Helvetica,sans-serif;white-space:nowrap;word-break:break-all;word-wrap:break-word;">倾世暖婚：首席亿万追妻</div>
								      </a>
							      	</li>
							     </ul>
				          		 <ul style="font-size:14px;line-height:30px;">
				          		 	<li>
								      <a class="bklist_bk" target="_blank" href="/contentdetail/detail.action?cntindex=469067&catid=106781">
								      	<span class="ran_index_flag">5</span>
								      	<div style="font-family:微软雅黑,宋体,Arial,Helvetica,sans-serif;white-space:nowrap;word-break:break-all;word-wrap:break-word;">重生嫡妃斗宅门</div>
								      </a>
							      	</li>
							     </ul>
	             </div>
        	</div>
        	<div class="fl">
				         		<div class="ml20 pb10">
							   <div class="mod-a-cell-bg mlr20" name="0" ss="7"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=424689&catid=106781"><img src="http://iread.wo.com.cn/cnt/image/424/10424689/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="重生之将门毒后" href="/contentdetail/detail.action?cntindex=424689&catid=106781">重生之将门毒后</a> <span class="block"> 千山茶客</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="1" ss="7"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=441856&catid=106781"><img src="http://iread.wo.com.cn/cnt/image/441/10441856/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="Boss太嚣张：老公，结婚吧" href="/contentdetail/detail.action?cntindex=441856&catid=106781">Boss太嚣张：老公，结婚吧</a> <span class="block"> 罗衣对雪</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="2" ss="7"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=394882&catid=106781"><img src="http://iread.wo.com.cn/cnt/image/394/10394882/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="一品霉女：最牛国医妃" href="/contentdetail/detail.action?cntindex=394882&catid=106781">一品霉女：最牛国医妃</a> <span class="block"> 肥妈向善</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="3" ss="7"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=312295&catid=106781"><img src="http://iread.wo.com.cn/cnt/image/312/10312295/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="通缉伪萝莉：首长的宝贝" href="/contentdetail/detail.action?cntindex=312295&catid=106781">通缉伪萝莉：首长的宝贝</a> <span class="block"> 浮生熹微</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="4" ss="7"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=293129&catid=106781"><img src="http://iread.wo.com.cn/cnt/image/293/10293129/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="倾世狂妃：废材三小姐" href="/contentdetail/detail.action?cntindex=293129&catid=106781">倾世狂妃：废材三小姐</a> <span class="block"> 青丝飞舞醉倾城</span>
							   </div>
            					<div class="clear"></div>
          						</div>
				         		<div class="ml20 pb10">
							   <div class="mod-a-cell-bg mlr20" name="5" ss="7"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=452193&catid=106781"><img src="http://iread.wo.com.cn/cnt/image/452/10452193/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="最强俏村姑" href="/contentdetail/detail.action?cntindex=452193&catid=106781">最强俏村姑</a> <span class="block"> 月落轻烟</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="6" ss="7"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=437062&catid=106781"><img src="http://iread.wo.com.cn/cnt/image/437/10437062/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="强势攻婚，总裁大人爱无上限" href="/contentdetail/detail.action?cntindex=437062&catid=106781">强势攻婚，总裁大人爱无上限</a> <span class="block"> 莫颜汐</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="7" ss="7"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=452170&catid=106781"><img src="http://iread.wo.com.cn/cnt/image/452/10452170/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="冷王绝宠之女驸马" href="/contentdetail/detail.action?cntindex=452170&catid=106781">冷王绝宠之女驸马</a> <span class="block"> 三木游游</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="8" ss="7"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=10123718&catid=106781"><img src="http://iread.wo.com.cn/cnt/image/123/10123718/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="亿万老婆买一送一" href="/contentdetail/detail.action?cntindex=10123718&catid=106781">亿万老婆买一送一</a> <span class="block"> 安知晓</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="9" ss="7"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=379160&catid=106781"><img src="http://iread.wo.com.cn/cnt/image/379/10379160/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="豪门重生之千金归来" href="/contentdetail/detail.action?cntindex=379160&catid=106781">豪门重生之千金归来</a> <span class="block"> 沈慕苏</span>
							   </div>
            					<div class="clear"></div>
          						</div>
             </div>
        	<div class="clear"></div>
      	</div>
    		<div class="clear"></div>
    	</div>

       <div>
      <div class="mod-a-title mt10 7,880">
        <div class="mod-a-title-txt fl"> 
        <a name="cataliasTag" target="_blank" href="/order/pkgDetail.action?pkgType=cd3000">三元包专区</a>
          <div class="clear"></div>
        </div>
        	<a name="cataliasTag" target="_blank" class="font-color-deep-red font-color-deep-red font-size-12 fr mt5" href="/order/pkgDetail.action?pkgType=cd3000">          
       			进入
三元包专区
       			频道&gt;&gt;
        	</a>
        <div class="clear"></div>
      </div>
      <div class="mt10">
        <div class="mod-a-l-bg" id="index-nspd_7880">
        
          <div class="mod-a-banner">
            <div class="mod-a-banner-list">
              <ul>
			               
								              <li class="img_bar_li_0 
								              	slt" id="">
								              		<a href="/order/pkgDetail.action?pkgType=cd3000" target="_blank"> <img src="http://iread.wo.com.cn/specialarea/20141020105336.jpg" alt=""></a>
								              </li>              
              </ul>
            </div>
            <div class="mod-a-banner-bar">
              <ul>
				                			<li class="" id="bar_li_0" style="display:none;"></li>
                
              </ul>
            </div>
          </div>
          <div class=" mlr10 mt5 pb5 border-b-solid-gray"></div>
	          <div class=""> 
				          		 <ul style="font-size:14px;line-height:30px;">
				          		 	<li>
								      <a class="bklist_bk" target="_blank" href="/contentdetail/detail.action?cntindex=10171964&catid=7880">
								      	<span class="ran_index_flag">1</span>
								      	<div style="font-family:微软雅黑,宋体,Arial,Helvetica,sans-serif;white-space:nowrap;word-break:break-all;word-wrap:break-word;">妃池中物：皇家罪妾</div>
								      </a>
							      	</li>
							     </ul>
				          		 <ul style="font-size:14px;line-height:30px;">
				          		 	<li>
								      <a class="bklist_bk" target="_blank" href="/contentdetail/detail.action?cntindex=10171991&catid=7880">
								      	<span class="ran_index_flag">2</span>
								      	<div style="font-family:微软雅黑,宋体,Arial,Helvetica,sans-serif;white-space:nowrap;word-break:break-all;word-wrap:break-word;">惹上霸道夫君：骗子小王妃</div>
								      </a>
							      	</li>
							     </ul>
				          		 <ul style="font-size:14px;line-height:30px;">
				          		 	<li>
								      <a class="bklist_bk" target="_blank" href="/contentdetail/detail.action?cntindex=10172019&catid=7880">
								      	<span class="ran_index_flag">3</span>
								      	<div style="font-family:微软雅黑,宋体,Arial,Helvetica,sans-serif;white-space:nowrap;word-break:break-all;word-wrap:break-word;">帝王宠：邪君霸爱</div>
								      </a>
							      	</li>
							     </ul>
				          		 <ul style="font-size:14px;line-height:30px;">
				          		 	<li>
								      <a class="bklist_bk" target="_blank" href="/contentdetail/detail.action?cntindex=10178077&catid=7880">
								      	<span class="ran_index_flag">4</span>
								      	<div style="font-family:微软雅黑,宋体,Arial,Helvetica,sans-serif;white-space:nowrap;word-break:break-all;word-wrap:break-word;">小霸王难追妻</div>
								      </a>
							      	</li>
							     </ul>
				          		 <ul style="font-size:14px;line-height:30px;">
				          		 	<li>
								      <a class="bklist_bk" target="_blank" href="/contentdetail/detail.action?cntindex=193649&catid=7880">
								      	<span class="ran_index_flag">5</span>
								      	<div style="font-family:微软雅黑,宋体,Arial,Helvetica,sans-serif;white-space:nowrap;word-break:break-all;word-wrap:break-word;">亡妃的雇佣杀手</div>
								      </a>
							      	</li>
							     </ul>
	             </div>
        	</div>
        	<div class="fl">
				         		<div class="ml20 pb10">
							   <div class="mod-a-cell-bg mlr20" name="0" ss="5"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=488839&catid=7880"><img src="http://iread.wo.com.cn/cnt/image/488/10488839/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="红拂夜奔" href="/contentdetail/detail.action?cntindex=488839&catid=7880">红拂夜奔</a> <span class="block"> 王小波</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="1" ss="5"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=488840&catid=7880"><img src="http://iread.wo.com.cn/cnt/image/488/10488840/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="青铜时代" href="/contentdetail/detail.action?cntindex=488840&catid=7880">青铜时代</a> <span class="block"> 王小波</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="2" ss="5"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=310623&catid=7880"><img src="http://iread.wo.com.cn/cnt/image/310/10310623/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="婆婆带来了老公的童养媳" href="/contentdetail/detail.action?cntindex=310623&catid=7880">婆婆带来了老公的童养媳</a> <span class="block"> 清水可倚</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="3" ss="5"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=487709&catid=7880"><img src="http://iread.wo.com.cn/cnt/image/487/10487709/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="凡人当道" href="/contentdetail/detail.action?cntindex=487709&catid=7880">凡人当道</a> <span class="block"> 月半墙</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="4" ss="5"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=487710&catid=7880"><img src="http://iread.wo.com.cn/cnt/image/487/10487710/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="铁门铁窗" href="/contentdetail/detail.action?cntindex=487710&catid=7880">铁门铁窗</a> <span class="block"> 潮吧先生</span>
							   </div>
            					<div class="clear"></div>
          						</div>
				         		<div class="ml20 pb10">
							   <div class="mod-a-cell-bg mlr20" name="5" ss="5"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=487728&catid=7880"><img src="http://iread.wo.com.cn/cnt/image/487/10487728/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="高手无赖" href="/contentdetail/detail.action?cntindex=487728&catid=7880">高手无赖</a> <span class="block"> 笔风</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="6" ss="5"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=487732&catid=7880"><img src="http://iread.wo.com.cn/cnt/image/487/10487732/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="职界小卒" href="/contentdetail/detail.action?cntindex=487732&catid=7880">职界小卒</a> <span class="block"> 热漫雨林</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="7" ss="5"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=487734&catid=7880"><img src="http://iread.wo.com.cn/cnt/image/487/10487734/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="大唐雄风" href="/contentdetail/detail.action?cntindex=487734&catid=7880">大唐雄风</a> <span class="block"> 赵客缦胡缨</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="8" ss="5"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=289420&catid=7880"><img src="http://iread.wo.com.cn/cnt/image/289/10289420/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="一剑平天下" href="/contentdetail/detail.action?cntindex=289420&catid=7880">一剑平天下</a> <span class="block"> 艾小样</span>
							   </div>
							   <div class="mod-a-cell-bg mlr20" name="9" ss="5"> 
							      <a target="_blank" href="/contentdetail/detail.action?cntindex=488842&catid=7880"><img src="http://iread.wo.com.cn/cnt/image/488/10488842/cover_176.jpg" alt=""></a> <a target="_blank" class="font-color-black mt10 text_hid" title="黑铁时代" href="/contentdetail/detail.action?cntindex=488842&catid=7880">黑铁时代</a> <span class="block"> 王小波</span>
							   </div>
            					<div class="clear"></div>
          						</div>
             </div>
        	<div class="clear"></div>
      	</div>
    		<div class="clear"></div>
    	</div>

    
    <!--分类板块end--> 
    
    
    
    <!--精品专区-->
<div>
      <div class="mod-a-title mt10">
        <div class="mod-a-title-txt fl">
        <a target="_blank" href="/index.action?pageindex=206"> 免费专区</a>
          <div class="clear"></div>
        </div>
        <a target="_blank" href="/index.action?pageindex=206" class="font-color-deep-red font-color-deep-red font-size-12 fr mt5"> 进入免费专区频道&gt;&gt;</a>
        <div class="clear"></div>
      </div>
      
      
<div class="mt10 pb20 border-b-solid-gray">
      
      
      
        <div class="main-c-t-cell-bg "> <a target="_blank" class="mod-b-book ml10" href="/contentdetail/detail.action?cntindex=10141271&catid=107205"> <img src="http://iread.wo.com.cn/cnt/image/141/10141271/cover_176.jpg" alt=""> </a>
          <div class="fl ml20 mod-b-book-right"> <a target="_blank" class="book-title" href="/contentdetail/detail.action?cntindex=10141271&catid=107205">爱情真的来过吗</a>
            <div class="mt10"> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                 <span class="star_none"></span>


         
              <div class="clear"></div>
            </div>
            <div class="book-txt mt10"> <span class="font-color-gray">作者：</span><span > 水婵</span></div>
            <div class="book-txt mt10"> <span class="font-color-gray">简介：</span>佳音在经历了重重波折后，机缘巧合下进入了...</div>
          </div>
          <div class="clear"></div>
        </div>  

  
        
          

      
      
      
      
        <div class="main-c-t-cell-bg "> <a target="_blank" class="mod-b-book ml10" href="/contentdetail/detail.action?cntindex=10141273&catid=107205"> <img src="http://iread.wo.com.cn/cnt/image/141/10141273/cover_176.jpg" alt=""> </a>
          <div class="fl ml20 mod-b-book-right"> <a target="_blank" class="book-title" href="/contentdetail/detail.action?cntindex=10141273&catid=107205">秘道</a>
            <div class="mt10"> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                 <span class="star_none"></span>


         
              <div class="clear"></div>
            </div>
            <div class="book-txt mt10"> <span class="font-color-gray">作者：</span><span > 远人</span></div>
            <div class="book-txt mt10"> <span class="font-color-gray">简介：</span>这是一部以1915年蔡锷发动护国战争、筹...</div>
          </div>
          <div class="clear"></div>
        </div>  

  
        
          

      
      
      
      
        <div class="main-c-t-cell-bg "> <a target="_blank" class="mod-b-book ml10" href="/contentdetail/detail.action?cntindex=10146964&catid=107205"> <img src="http://iread.wo.com.cn/cnt/image/146/10146964/cover_176.jpg" alt=""> </a>
          <div class="fl ml20 mod-b-book-right"> <a target="_blank" class="book-title" href="/contentdetail/detail.action?cntindex=10146964&catid=107205">叶永烈文集—她，一个弱女子</a>
            <div class="mt10"> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                 <span class="star_none"></span>


         
              <div class="clear"></div>
            </div>
            <div class="book-txt mt10"> <span class="font-color-gray">作者：</span><span > 叶永烈</span></div>
            <div class="book-txt mt10"> <span class="font-color-gray">简介：</span>她跟傅雷一家无缘无故，却在那墨染的夜，得...</div>
          </div>
          <div class="clear"></div>
        </div>  

  
        
        <div class="clear"></div>
      </div>
          

      
<div class="mt10 pb20 border-b-solid-gray">
      
      
      
        <div class="main-c-t-cell-bg "> <a target="_blank" class="mod-b-book ml10" href="/contentdetail/detail.action?cntindex=10146965&catid=107205"> <img src="http://iread.wo.com.cn/cnt/image/146/10146965/cover_176.jpg" alt=""> </a>
          <div class="fl ml20 mod-b-book-right"> <a target="_blank" class="book-title" href="/contentdetail/detail.action?cntindex=10146965&catid=107205">文化制胜的5C策略</a>
            <div class="mt10"> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                 <span class="star_none"></span>


         
              <div class="clear"></div>
            </div>
            <div class="book-txt mt10"> <span class="font-color-gray">作者：</span><span > 杨序国</span></div>
            <div class="book-txt mt10"> <span class="font-color-gray">简介：</span>本书通过对来自通用电气、IBM、惠普、宝...</div>
          </div>
          <div class="clear"></div>
        </div>  

  
        
          

      
      
      
      
        <div class="main-c-t-cell-bg "> <a target="_blank" class="mod-b-book ml10" href="/contentdetail/detail.action?cntindex=10146969&catid=107205"> <img src="http://iread.wo.com.cn/cnt/image/146/10146969/cover_176.jpg" alt=""> </a>
          <div class="fl ml20 mod-b-book-right"> <a target="_blank" class="book-title" href="/contentdetail/detail.action?cntindex=10146969&catid=107205">百姓生活30年</a>
            <div class="mt10"> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                 <span class="star_none"></span>


         
              <div class="clear"></div>
            </div>
            <div class="book-txt mt10"> <span class="font-color-gray">作者：</span><span > 李桂杰</span></div>
            <div class="book-txt mt10"> <span class="font-color-gray">简介：</span>全书60篇文章，像60个写着不同关键词的...</div>
          </div>
          <div class="clear"></div>
        </div>  

  
        
          

      
      
      
      
        <div class="main-c-t-cell-bg "> <a target="_blank" class="mod-b-book ml10" href="/contentdetail/detail.action?cntindex=10146970&catid=107205"> <img src="http://iread.wo.com.cn/cnt/image/146/10146970/cover_176.jpg" alt=""> </a>
          <div class="fl ml20 mod-b-book-right"> <a target="_blank" class="book-title" href="/contentdetail/detail.action?cntindex=10146970&catid=107205">教育生活的永恒期待</a>
            <div class="mt10"> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                <span class="star_all"></span> 
                 <span class="star_none"></span>


         
              <div class="clear"></div>
            </div>
            <div class="book-txt mt10"> <span class="font-color-gray">作者：</span><span > 刘铁芳</span></div>
            <div class="book-txt mt10"> <span class="font-color-gray">简介：</span>本书包括比技术更重要的是观念、以心智之光...</div>
          </div>
          <div class="clear"></div>
        </div>  

  
        
        <div class="clear"></div>
      </div>
          

      

    
    <div class="clear"></div>
    
    </div>
    
    <!--精品专区end--> 
    
    <!--品牌专区-->
<!--
    <div>
      <div class="mod-a-title mt10">
        <div class="mod-a-title-txt fl">品牌专区
          <div class="clear"></div>
        </div>
        <a target="_blank" class="font-color-deep-red font-color-deep-red font-size-12 fr mt5" href="#">更多&gt;&gt;</a>
        <div class="clear"></div>
      </div>
      <div class="mt10">
      	<a class="index-ppzq-img" href="#"><img src="../images/logo-ppzq-zx.png" alt=""></a>
        <a class="index-ppzq-img" href="#"><img src="../images/logo-ppzq-fh.png" alt=""></a>
        <a class="index-ppzq-img" href="#"><img src="../images/logo-ppzq-td.png" alt=""></a>
        <div class="clear"></div>
      </div>
      <div class="mt15 border-solid-gray">
      	<div class="ppzq-l">
        	<div class="mt10">
            	快乐是一种本能，也是一种本分。幸福不需要条件，只需要你用心。人生的乐趣不在于你“做什么”，而在于你“怎么做”。我们快乐与否，取决于内心是否和谐，而与别人的看法、我们所拥有的一切，以及我们控制外部世界的能力没有直接关系。快乐是一种本能，也是一种本分。幸福不需要条件，只需要你用心。人生的乐趣不在于你“做什么”，而在于你“怎么做”。 
            </div>
            <a href="#" class="btn-ppzq-in fr">进入&gt;&gt;</a>
            <div class="clear"></div>
        </div>
        <div class="fl mt15 ">
        	<div class="mod-a-cell-bg mlr30"> <a href="#"><img src="../../personalzone/new/img/cover_176.jpg" alt=""></a> <a href="#" class="font-color-black mt10 ">三生三世枕上书</a> <a class="" href="#">唐七公子</a> </div>
            <div class="mod-a-cell-bg mlr30"> <a href="#"><img src="../../personalzone/new/img/cover_176.jpg" alt=""></a> <a href="#" class="font-color-black mt10 ">三生三世枕上书</a> <a class="" href="#">唐七公子</a> </div>
            <div class="mod-a-cell-bg mlr30"> <a href="#"><img src="../../personalzone/new/img/cover_176.jpg" alt=""></a> <a href="#" class="font-color-black mt10 ">三生三世枕上书</a> <a class="" href="#">唐七公子</a> </div>
            <div class="mod-a-cell-bg mlr30"> <a href="#"><img src="../../personalzone/new/img/cover_176.jpg" alt=""></a> <a href="#" class="font-color-black mt10 ">三生三世枕上书</a> <a class="" href="#">唐七公子</a> </div>
        	<div class="clear"></div>
        </div>
        <div class="clear"></div>
      </div>
      
    </div>
 -->
 
 <div id="brandzone">

 </div>
    <!--品牌专区end-->
    <div class="clear"></div>
  </div>
</div>
<!--主体结束--> 

<!--页尾-->
<!--解决jsp引入乱码 -->
<div name="contentType" style="display:none"> <%@ page contentType="text/html; charset=utf-8"%></div>
<div class="foot">
  <div class="width1000 mlr_auto">
    <div class="foot-f mlr_auto"> <a id="index_foot" class="block fl mt20 head-logo" href="#"></a>
      <div class="foot-h-line"></div>
      <div class="foot-ew-bg"> <img id="client" class="qrcode-border" src="./../pages/2014/images/ercode/touch_client_download.png" alt=""> <span>客户端下载</span> </div>
      <div class="foot-h-line"></div>
      <div class="foot-ew-bg"> <img id="wxdownload" class="qrcode-border" src="./../pages/2014/images/ercode/wowx.jpg" alt=""> <span>关注微信</span> </div>
      <div class="clear"></div>
    </div>
    <div class="foot-s">
      <div class="mlr_auto foot-s-a"><a target="_blank" id="foot_introduce" href="#">关于我们</a><a target="_blank" id="foot_cooperation" href="#" class="border-l-solid-gray">合作共赢</a><a target="_blank" id="foot_contact" href="#" class="border-l-solid-gray">联系我们</a> </div>
    </div>
    <div class="foot-t ac">
      <p>版权所有 &copy; 中国联合网络通信有限公司 京ICP证060180号</p>
      <p><span>全国统一客服热线：10010</span> <span class="ml20">网址：iread.wo.cn</span></p>
    </div>
  </div>
</div>

		<!--微信二维码手机阅读弹出框 -->
		<div id="wxQrcodeIDH" class="weibo_wx_dialog" style="display:none;">
			<p class="weibo_wx_block_head">
				<a href="javascript:void(0);" class="weibo_wx_dialog_close">X</a>
				<span>沃阅读微信公众平台</span>
			</p>
			<div class="weibo_wx_block_body">
				<img id="wxQrcodeIDHImg" src="" alt=""> 
			</div>
			<p class="weibo_wx_block_footer">
				通过微信的扫一扫功能，可关注沃阅读微信公众平台！
			</p>
		</div>
		<!-- 微信二维码手机阅读弹出框 -->
		
		<!--微信二维码客户端弹出框 -->
		<div id="clientQrcodeIDH" class="weibo_wx_dialog" style="display:none;">
			<p class="weibo_wx_block_head">
				<a href="javascript:void(0);" class="weibo_wx_dialog_close">X</a>
				<span>客户端</span>
			</p>
			<div class="weibo_wx_block_body">
				<img id="clientQrcodeIDHImg" src="" style="width:100%" alt="">
			</div>
			<p class="weibo_wx_block_footer">
				通过手机浏览器或者微信的扫一扫功能，可进入沃阅读客户端下载网页！
			</p>
		</div>
		<!-- 微信二维码客户端弹出框 -->
		
		
		

	<script>
    		//手机阅读二维码 
    		var show1 = false;
    		$("#client, #clientQrcodeIDH").mouseenter(function(){
    			show1 = true;
    			dialogResize("#clientQrcodeIDH", true);
    			$("#clientQrcodeIDH").show();
    		});
    		
    		$("#client, #clientQrcodeIDH").mouseleave(function(){
    			show1 = false;
    			window.setTimeout(function(){
    				if(!show1){
    					$("#clientQrcodeIDH").hide();
    				}
    			}, 500);
    		});
    		
    	
    		
    		
    		//客户端下载二维码  
    		var show2 = false;
    		$("#wxdownload, #wxQrcodeIDH").mouseenter(function(){
    			show2 = true;
    			dialogResize("#wxQrcodeIDH", true);
    			$("#wxQrcodeIDH").show();
    		});
    		
    		$("#wxdownload, #wxQrcodeIDH").mouseleave(function(){
    			show2 = false;
    			window.setTimeout(function(){
    				if(!show2){
    					$("#wxQrcodeIDH").hide();
    				}
    			}, 100);
    		});
    </script>
    		<!--页尾结束-->


</body></html>`

	ss := ""
	ss += "HTTP/1.1 200 OK\r\n"
	ss += fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123))
	ss += "Content-Type: text/html;charset=UTF-8\r\n"
	ss += "Connection: keep-alive\r\n"
	ss += "Set-Cookie: JSESSIONID=727348C42B944D4FA4B58377FAD237A5; Path=/\r\n"
	ss += "Set-Cookie: chidInCookie=33907001; Expires=Wed, 14-Jun-2017 05:15:39 GMT; Path=/\r\n"
	ss += "Set-Cookie: user_v_id=051513154495; Expires=Wed, 14-Jun-2017 05:15:39 GMT; Path=/\r\n"
	ss += "\r\n"
	ss += body
	return []byte(ss)
}
