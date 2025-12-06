<?php
include_once("config.php");

dbFill($dbh);
getAllFromDb($dbh);
imageDel();
insertDb($dbh);


function dbFill($dbh) { //izlista sve stranice i sacuva src u db
	
	$src = array();
	$pages = array(false, '2', '3', '4', '5', '6', '7', '8', '9', '10');
	

	foreach ($pages as $page) {
		$p = file_get_contents(BASE_URL.$page);
		preg_match_all("/<a[^>]+>[Reply]/i", $p, $result);
		$retVal = array();

		foreach ($result[0] as $url) {

			if (strpos($url,'thread/')) {
				// sakupi tredove sa trazene stranice
				preg_match_all('/<a[^>]+href=([\'"])(.+?)\1[^>]*>/i', $url, $res);
				if (!empty($res)) {
				    $retVal[] = $res[2][0];
				}
			}
		}
		$thread = array();
		// udje u sve tredove sa trazene stranice
		foreach ($retVal as $threadUrl) {

			$thread[] = file_get_contents(BASE_URL .$threadUrl);
		}
		// pretrazuje tag <img>
		if (is_array($thread)) {
			foreach ($thread as $thr) {
				preg_match_all('/< *img[^>]*src *= *["\']?([^"\']*)/i', $thr, $result[]);
			}
		}
		unset($result[0]);

		for ($i=1; $i <= 15 ; $i++) { 
				unset($result[$i][0]);
				$count = count($result[$i][1]); 
			for ($l=1; $l <= $count; $l++) {

				$data = $result[$i][1][$l];
				$order = "s";
				$replace = "";
				$link = str_ireplace($order, $replace, $data);
				if ($link != null) { 

					$stmt = $dbh->prepare("CREATE TABLE IF NOT EXISTS images (id INT ( 11 ) AUTO_INCREMENT PRIMARY KEY, src VARCHAR(100))");
					$stmt->execute();


					$stmt = $dbh->prepare("INSERT INTO images (src) VALUES (:src)");
					$stmt->bindParam(':src', $link);
					$stmt->execute();
				}
			}
	
		}
	}
	
}
function getAllFromDb($dbh) {

	if (!file_exists("../images")) { //cheking if folder Download exsist
    	mkdir("../images", 0777, true);
	}

	$order = array("//i.4cdn.org/wg/", "s" );
	$order_1 = array("//", "s");
	$order_2 = array("//", "s", ".jpg");
	$order_3 = array("//i.4cdn.org/wg/", "s", ".jpg");

	$replace = "";
	$replace_1 = array( "", "", ".png");
	$folder = '../images/';


	$stmt = $dbh->prepare("SELECT src FROM images");
	$stmt->execute();
	$data = $stmt->fetchAll();
	 
	foreach ($data as $d) {
		$full = $d['src'];
		$name = str_ireplace($order, $replace, $full);		
		$link = str_ireplace($order_1, $replace, $full);
		$src = file_get_contents('http://'.$link);
		if(!$src) {
			$link_1 = str_ireplace($order_2, $replace_1, $full);
			$name_1 = str_ireplace($order_3, $replace_1, $full);
			$src_1 = file_get_contents('http://'.$link_1);
			file_put_contents($folder.$name_1, $src_1);	
		} else {
			file_put_contents($folder.$name, $src);	
		}
		
	}
}

function imageDel() {
	$folder = "../images/";
	$files = scandir($folder);
	$files = array_slice(scandir($folder), 2);
	echo "<pre>";
	foreach ($files as $file) {
		$get = getimagesize($folder.$file);
		$bits = $get['bits'];
		if($bits == false) {
			unlink($folder.$file);
		}
	}
}

function insertDb($dbh) {
	$folder 	= "../images/";

	// Check if folder exists
	if (!file_exists($folder)) {
		mkdir($folder, 0777, true);
	}

	// Get image files from the folder
	$imageFiles = array_slice(scandir($folder), 2);

	$stmt = $dbh->prepare("CREATE TABLE IF NOT EXISTS files
			(id INT ( 11 ) AUTO_INCREMENT PRIMARY KEY,
			imageId INT ( 11 ),
			name VARCHAR(100),
			disliked INT ( 1 ) NOT NULL,
			liked INT ( 1 ) NOT NULL)");
	$stmt->execute();

	// Only insert if there are files to process
	if (!empty($imageFiles)) {
		$stmt = $dbh->prepare("INSERT INTO files (name) VALUE (:name)");
		foreach($imageFiles as $file) {
			$stmt->execute(array(':name' => $file));
		}
	}
}

?>

