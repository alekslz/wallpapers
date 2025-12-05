$( document ).ready(function() {

	$.ajax({
		type: "GET",
		url: "/api/count",
		success: function(data) {
			var track_load = 0;
			var loading = false;
			var total_groups = data;
			//alert(track_load);

			$(".images").load("/api/images", {'group_no':track_load}, function() {track_load++;});
			$(window).scroll(function() { //detect page scroll
        	console.log(track_load);
		        if($(window).scrollTop() + $(window).height() == $(document).height())  //user scrolled to bottom of the page?
		        {

		            if(track_load <= total_groups && loading==false) //there's more data to load
		            {
		                loading = true; //prevent further ajax loading
		                //$('.animation_image').show(); //show loading image

		                //load data from the server using a HTTP POST request
		                $.post("/api/images",{'group_no': track_load}, function(data){
		                                    
		                    $(".images").append(data); //append received data into the element
		                    //hide loading image
		                    //$('.animation_image').hide(); //hide loading image once data is received
		                    
		                    track_load++; //loaded group increment
		                    loading = false; 
		                
		                }).fail(function(xhr, ajaxOptions, thrownError) { //any errors?
		                    
		                    alert(thrownError); //alert with HTTP error
		                    //$('.animation_image').hide(); //hide loading image
		                    loading = false;
		                
		                });
		                
		            }
		        }
		    });

		} 
	})

});
