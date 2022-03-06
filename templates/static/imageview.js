// Handle Keypresses
$(document).on("keydown", function (e) {
	switch (e.which) {
		case 37: // left
			window.location.href = $("#previous").attr("href");
			break;
		case 39: // right
			window.location.href = $("#next").attr("href");
			break;
		case 70: // f
			let media = $("#media");
			switch (media.attr("class")) {
				case undefined:
				case "":
					media.addClass("bigscreen");
					break;
				case "bigscreen":
					media.removeClass("bigscreen");
					media.addClass("biggerscreen");
					break;
				case "biggerscreen":
					media.removeClass("biggerscreen");
					media.addClass("fullscreen");
					break;
				case "fullscreen":
					media.removeClass("fullscreen");
					break;
			}
			break;
	}
});

// Do timer
const urlParams = new URLSearchParams(window.location.search);
if (urlParams.get("timer") != null) {
	setTimeout(function () {
		window.location.href = $("#next").attr("href");
	}, urlParams.get("timer") * 1000);
}

// Resize svg to image size
let image = $("#media");
let svg = $("#faceCanvas");
svg.css("top", image.css("top"));
svg.css("width", image.width() + "px");
svg.css("height", image.height() + "px");
window.addEventListener("resize", function (event) {
	svg.css("top", image.css("top"));
	svg.css("width", image.width() + "px");
	svg.css("height", image.height() + "px");
});

// Setup context menu for face boxes
$(".facebox").on("contextmenu", function (e) {
	var top = e.pageY - 10;
	var left = e.pageX - 90;
	$("#context-menu")
		.attr("data", $(e.target).attr("data"))
		.css({
			display: "absolute",
			top: top,
			left: left,
		})
		.show();
	return false;
});

$("#context-menu #viewFace").on("click", function (e) {
	let menu = $(e.target).parent();
	window.location.href =
		"/face" + menu.attr("image-url") + "?face=" + menu.attr("data");
	return false;
});

$("#context-menu #searchFace").on("click", function (e) {
	let menu = $(e.target).parent();
	window.location.href =
		"/api/aws/recognize" + menu.attr("image-url") + "?face=" + menu.attr("data");
	return false;
});

$("#context-menu a").on("click", function () {
	$(this).parent().hide();
});

$(document)
	.on("click", function () {
		$("#context-menu").hide();
	})
	.on("contextmenu", function (e) {
		$("#context-menu").hide();
	});
