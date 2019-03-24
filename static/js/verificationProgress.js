top.progressStart = 0;
top.progressStartTime = 0;

var durationFormatter = function( seconds ) {


    var interval = Math.floor(seconds / 31536000);

    if (interval > 1) {
        return interval + " years";
    }
    interval = Math.floor(seconds / 2592000);
    if (interval > 1) {
        return interval + " months";
    }
    interval = Math.floor(seconds / 86400);
    if (interval > 1) {
        return interval + " days";
    }
    interval = Math.floor(seconds / 3600);
    if (interval > 1) {
        return interval + " hours";
    }
    interval = Math.floor(seconds / 60);
    if (interval > 1) {
        return interval + " minutes";
    }
    return Math.floor(seconds) + " seconds";

}

var loadVerificationProgress = function() {
    console.log( "loading verification progress" );

    var request = new XMLHttpRequest();
    request.open('GET', '/verificationprogress');
    request.responseType = 'text';

    request.onerror = function(e) {  }

    request.onload = function() {
        if (request.readyState==4 && request.status==200){
            var result;

            try {
                result = JSON.parse( request.response )
            } catch( e ) {
                console.log( e );
            }

            if( result && result.verificationprogress ) {
                result.verificationprogress = parseInt(result.verificationprogress*100000)/100000;
                if( top.progressStart === 0 ) {
                    top.progressStart = result.verificationprogress;
                    top.progressStartTime = (new Date())/1000;
                } else {
                    var deltaS = (new Date())/1000 - top.progressStartTime;
                    var deltaP = result.verificationprogress - top.progressStart;

                    var eta = deltaS/deltaP;
                    var pBar = document.getElementById('progress-bar');
                    if( pBar ) {
                        pBar.style.width = (result.verificationprogress*100)+"%";
                    }
                    var pText = document.getElementById('progress-text');
                    if( pText ) {
                        pText.innerText = "Verification complete in "+durationFormatter(eta);
                    }
                }
            }
        }
    };
    request.send();
}

loadVerificationProgress();
setInterval( loadVerificationProgress, 5000 );
