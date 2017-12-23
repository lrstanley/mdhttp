// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package mdhttp

const htmlTemplate = `<!DOCTYPE html>
<html>
	<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<meta charset="UTF-8">
	<title>{{.Title}} &middot; {{ .FileInfo.Name }}</title>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/2.9.0/github-markdown.min.css" integrity="sha256-AcEoN20lbehGx2iA6fVk4geoblRBOAUmz+gEM9QE59Y=" crossorigin="anonymous" />
	<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/2.2.4/jquery.min.js" integrity="sha256-BbhdlvQf/xTY9gja0Dq3HiwQF8LaCRTXxZKRutelT44=" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css" integrity="sha256-eZrrJcwDc/3uDhsdt61sL2oOBY362qM3lon1gyExkL0=" crossorigin="anonymous" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.12.0/styles/default.min.css" integrity="sha256-Zd1icfZ72UBmsId/mUcagrmN7IN5Qkrvh75ICHIQVTk=" crossorigin="anonymous" />
	<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.12.0/highlight.min.js" integrity="sha256-/BfiIkHlHoVihZdc6TFuj7MmJ0TWcWsMXkeDFwhi0zw=" crossorigin="anonymous"></script>
	<script>hljs.initHighlightingOnLoad();</script>
	<style>
		html, body {
			font-family: -apple-system,BlinkMacSystemFont,"Segoe UI",Helvetica,Arial,sans-serif,"Apple Color Emoji","Segoe UI Emoji","Segoe UI Symbol";
		}

		.main {
			min-width: 200px;
			max-width: 980px;
			margin: 40px auto;
		}

		.markdown-title {
			display: block;
			padding: 9px 10px 10px;
			margin: 0;
			font-size: 14px;
			line-height: 17px;
			background-color: #f6f8fa;
			border: 1px solid rgba(27,31,35,0.15);
			border-bottom: 0;
			border-radius: 3px 3px 0 0;
		}

		.markdown-title > span { margin-right: 10px; }

		.markdown-body {
			box-sizing: border-box;
			padding: 20px 45px;
			word-wrap: break-word;
			background-color: #fff;
			border: 1px solid #ddd;
			border-bottom-right-radius: 3px;
			border-bottom-left-radius: 3px;
		}

		@media (max-width: 767px) {
			.markdown-body { padding: 15px; }
		}

		h1 > .anchor, h2 > .anchor, h3 > .anchor, h4 > .anchor {
			font-size: 20px;
			position: relative;
			text-decoration: none;
			color: #639fe2;
		}

		h1 > .anchor:hover, h2 > .anchor:hover, h3 > .anchor:hover, h4 > .anchor:hover { color: #0366d6; }
		h1 > .anchor { top: 10px; }
		h2 > .anchor { top: 5px; }
		h3 > .anchor { top: 2.5px; }
		h4 > .anchor { font-size: 18px; }

		.footer {
			font-size: 13px;
			color: grey;
			margin: 10px 0;
		}
	</style>
</head>

<body>
	<div class="main">
		<h3 class="markdown-title">
			<span><i class="fa fa-book"></i> {{.Title}}</span>
			<span><i class="fa fa-file-code-o"></i> {{ .FileInfo.Name }}</span>
			<span><i class="fa fa-clock-o"></i> {{ .FileInfo.ModTime.Format "Mon, 02 Jan 2006 15:04:05 MST" }}</span>
		</h3>
		<div class="markdown-body">{{.HTMLTemplate}}</div>

		<footer class="footer">
			Generated with <a target="_blank" href="https://github.com/russross/blackfriday">blackfriday</a> &middot;
			Served with <a target="_blank" href="https://github.com/lrstanley/mdhttp">mdhttp</a>
	</div>

	<script type="text/javascript">
	$(function () {
		$.each($('.markdown-body h1, .markdown-body h2, .markdown-body h3, .markdown-body h4'), function(index, value) {
			var anchor = makeAnchor($(this).text());
			$(this).attr("id", anchor);

			$(this).html("<a class='anchor' href='#" + anchor + "'>#</a> " + $(this).html());
		});

		setTimeout(function() {
			if (window.location.hash != "") {
				el = document.getElementById(window.location.hash.replace("#", ""));
				if (el !== null) {
					el.scrollIntoView();
				}
			}
		}, 500);
	});

	var anchorArray = [];

	function makeAnchor(text) {
		text = text.toLowerCase().replace(/[^a-z0-9-_ ]/g, "").trim().replace(/[-_ ]+/g, "-");

		var count = 0;
		for (i = 0; i < anchorArray.length; i++) {
			if (text == anchorArray[i]) {
				count++
			}
		}

		anchorArray.push(text);

		if (count == 0) { return text; }

		return text + "-" + (count + 1);
	}
	</script>
</body>
</html>`
