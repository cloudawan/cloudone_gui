var moduleTimeSeriesChart = (function(){

	function selectColor(index) {
		switch(index) {
			case 0:
				return "red";
			case 1:
				return "blue";
			case 2:
				return "green";
			case 3:
				return "fuchsia";
			case 4:
				return "maroon";
			case 5:
				return "yellow";
			case 6:
				return "lime";
			case 7:
				return "purple";
			case 8:
				return "olive";
			case 9:
				return "teal";
			case 10:
				return "aqua";
			case 11:
				return "navy";
			case 12:
				return "brown";
			case 13:
				return "chocolate";
			case 14:
				return "hotpink";
			case 15:
				return "firebrick";
			case 16:
				return "lightseagreen";
			case 17:
				return "orangered";
			case 18:
				return "rosybrown";
			case 19:
				return "orchid";
			default:
				var value = index % 255;
				return "rgb(" + value + "," + value + "," + value + ")";
		}
	}

	var draw = function (divID, metadata, data) {
		// Example
		//var divID = "#idChartCpuUsageTotal";
		//var metadata = {
		//	title: "Title",
		//	lineName: ["line1", "line2", "line3"],
		//}
		//var data = [
		//	{y: [4,2,3], x: "2015-07-24 00:43:29"},
		//	{y: [4,2,3], x: "2015-07-24 00:43:30"},
		//	{y: [1,2,3], x: "2015-07-24 00:43:31"},
		//	{y: [4,2,3], x: "2015-07-24 00:43:32"},
		//	{y: [1,2,3], x: "2015-07-24 00:43:33"},
		//	{y: [4,2,3], x: "2015-07-24 00:43:34"},
		//	{y: [1,2,3], x: "2015-07-24 00:43:35"},
		//	{y: [1,2,3], x: "2015-07-24 00:43:36"},
		//	{y: [1,2,3], x: "2015-07-24 00:43:37"},
		//	{y: [1,2,3], x: "2015-07-24 00:43:38"},
		//]; 
		
		var parseDate = d3.time.format("%Y-%m-%d %H:%M:%S").parse;
		
		// Calculate the minimum and maximum of chart
		var xMin = parseDate(data[0].x);
		var xMax = parseDate(data[data.length-1].x);
		var yMin = Number.MAX_VALUE;
		var yMax = Number.MIN_VALUE;
		for (var i=0;i<data.length;i++) {
			for (var j=0;j<data[i].y.length;j++) {
				var value = data[i].y[j];
				if (value < yMin) {
					yMin = value;
				}
				if (value > yMax) {
					yMax = value;
				}
			}
		}
		yMax = Math.ceil(yMax * 1.2);
		
		var xLabelAmount = 10;
		var yLabelAmount = 10;
	
		var w = 1000;
		var h = 500;
		var m = 70;
		var x = d3.time.scale().domain([ xMin, xMax]).range([0, w]);
		var y = d3.scale.linear().domain([ yMin, yMax]).range([h, 0]);
	
		// Clean the previous chart
		$(divID).empty();
	
		var vis = d3.select(divID)
			.data([data])
			.append("svg:svg")
			.attr("width", w + m * 2)
			.attr("height", h + m * 2)
			.append("svg:g")
			.attr("transform", "translate(" + m + "," + m + ")");

		var xAxis = d3.svg.axis()
			.scale(x)
			.orient("bottom");

		vis.append("svg:g")
			.attr("transform", "translate(0," + h + ")")
			.attr("class", "axisX")
			.attr("height", 1)
			.call(xAxis);

		var rules = vis.selectAll("g.rule")
			.data(x.ticks(xLabelAmount))
			.enter().append("svg:g");
		
		rules.append("svg:text")
			.data(y.ticks(yLabelAmount))
			.attr("y", y)
			.attr("x", -10)
			.attr("dy", ".35em")
			.attr("text-anchor", "end")
			.text(y.tickFormat(",.0f"));

		for (var i=0; i < data[0].y.length; i++) {
			vis.append("svg:path")
				.attr("fill", "none")
				.attr("stroke", selectColor(i))
				.attr("stroke-width", 2)
				.attr("d", d3.svg.line()
					.x(function(d) { return x(parseDate(d.x)); })
					.y(function(d) { return y(d.y[i]); }));
			
			// Legend
			vis.append("svg:rect")
				.attr("x", w/4 - m/2)
				.attr("y", m/2 + 12 * (i+1))
				.attr("stroke", selectColor(i))
				.attr("stroke-width", 1)
				.attr("height", 1)
				.attr("width", 40);
				
			vis.append("svg:text")
				.attr("x", w/4 + m/2)
				.attr("y",(m/2 + 5) + 12 * (i+1))
				.text(metadata.lineName[i]);
		}
		
		// Add Title
		vis.append("svg:text")
			.attr("x", w/2 - m)
			.attr("y", m/2)
			.text(metadata.title);
	}

	return {
		draw: draw,
	};
			
})();