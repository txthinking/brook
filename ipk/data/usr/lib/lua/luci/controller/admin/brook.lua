module("luci.controller.admin.brook", package.seeall)
function index()
	entry({"admin", "brook"}, template("brook"),"Brook", 99).index=true
end
