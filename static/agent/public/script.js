$(document).ready(function () {
  $.ajax({
    url: METRICS_URL,
    success: function (result) {
      for (const checkId in result.checks) {
        const check = result.checks[checkId];
        const checkDiv = $("#check--"+checkId);
        const checkBarDiv = checkDiv.find(".main-check-bars");
        if (!check.operational) {
          checkDiv.addClass("main-check-failed");
        } else {
          checkDiv.addClass("main-check-success");
        }
        checkDiv.removeClass("main-check-loading");

        for (const metric of check.metrics) {
          const elem = $("<div></div>");
          elem.addClass("main-check-bar");
          if (metric.timeout) {
            elem.addClass("main-check-bar-timeout");
          } else if (metric.successful) {
            elem.addClass("main-check-bar-success");
          } else {
            elem.addClass("main-check-bar-failed");
          }
          const startTime = (new Date(metric.start_time)).toUTCString();
          const duration = (metric.duration / 1_000_000_000).toFixed(3);
          elem.attr("title", "Start Time: "+startTime+"\nDuration: "+duration+"s");
          checkBarDiv.append(elem);
        }
      }

      const operationalDiv = $(".main-operational-status");
      const failed = result.checks_down;
      if (failed === 0) {
        operationalDiv.addClass("main-operational-status-success");
      } else {
        operationalDiv.addClass("main-operational-status-failed");
        let text = "1 system down";
        if (failed > 1) {
          text = failed+" systems down";
        }
        $(".main-operational-status-failed-msg .main-operational-status-msg-text").text(text);
      }
      operationalDiv.removeClass("main-operational-status-loading");
    },
    error: function (err) {
      alert("Error: "+err.responseJSON.error);
    }
  });
});
