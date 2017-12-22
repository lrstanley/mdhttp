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
	<title>{{.Title}}</title>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/2.9.0/github-markdown.min.css" integrity="sha256-AcEoN20lbehGx2iA6fVk4geoblRBOAUmz+gEM9QE59Y=" crossorigin="anonymous" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.12.0/styles/default.min.css" integrity="sha256-Zd1icfZ72UBmsId/mUcagrmN7IN5Qkrvh75ICHIQVTk=" crossorigin="anonymous" />
	<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.12.0/highlight.min.js" integrity="sha256-/BfiIkHlHoVihZdc6TFuj7MmJ0TWcWsMXkeDFwhi0zw=" crossorigin="anonymous"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/2.2.4/jquery.min.js" integrity="sha256-BbhdlvQf/xTY9gja0Dq3HiwQF8LaCRTXxZKRutelT44=" crossorigin="anonymous"></script>
	<script>hljs.initHighlightingOnLoad();</script>
	<style>
		.markdown-body {
			box-sizing: border-box;
			min-width: 200px;
			max-width: 980px;
			margin: 0 auto;
			padding: 45px;
		}

		@media (max-width: 767px) {
			.markdown-body {
				padding: 15px;
			}
		}

		h2 > .anchor, h3 > .anchor, h4 > .anchor {
			font-size: 20px;
			position: relative;
			top: 5px;
			text-decoration: none;
			color: #639fe2;
		}

		h2 > .anchor:hover, h3 > .anchor:hover, h4 > .anchor:hover {
			color: #0366d6;
		}

		h3 > .anchor {
			top: 2.5px;
		}

		h4 > .anchor {
			font-size: 18px;
		}
	</style>
</head>

<body>
	<div class="main markdown-body">{{.HTMLTemplate}}</div>

	<script type="text/javascript">
	$(function () {
		$.each($('.markdown-body h2, .markdown-body h3, .markdown-body h4'), function(index, value) {
			var anchor = makeAnchor($(this).text());
			$(this).attr("id", anchor);

			$(this).html("<a class='anchor' href='#" + anchor + "'>#</a> " + $(this).html());
		});
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
