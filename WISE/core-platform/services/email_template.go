package services

const BusinessStatusChange = `
%s,<br/>
 %s has been %s.<br/>
Thanks,<br/>
Wise Customer Success Team`

// BusinessReviewEmailSubject  ..
const BusinessReviewEmailSubject = "Thank you for submitting your Wise Banking application!"

// BusinessReviewEmail business review status copy
const BusinessReviewEmail = `%s,<br/>
<br/>
Thank you for submitting your Wise Banking application! You are on your way to a better, faster banking experience.<br/>
<br/>
We are now reviewing your application and will get back in touch shortly. In the meantime, if you have additional questions, we are here to help.<br/>
<br/>
Contact the Wise Customer Success Team via hub.wise.us or send an email to chat@wise.us.<br/>
<br/>
We look forward to working with you!<br/>
Wise Customer Success Team<br/>
<br/>
© 2020 Wise Company.<br/>
Wise Company banking services provided by BBVA USA.<br/>
The Wise Company Visa Card is issued pursuant to license from Visa U.S.A.<br/>
<a href="https://hub.wise.us/en/collections/1642414-legal-disclosures-terms-and-policies"> Legal Disclosures, Terms and Policies </a><br/>
Need help? Email us at chat@wise.us.<br/>
`

// BusinessAddressBlock is a template for address block
const BusinessAddressBlock = `%s<br />
%s
%s, %s. %s
`

// BusinessApprovedEmailSubject  ..
const BusinessApprovedEmailSubject = ` Congratulations! Your Wise account is approved!! Let’s get you set up. `

// BusinessApprovedEmail  business status approved copy
const BusinessApprovedEmail = `Welcome to Wise, %s
<br /><br />
Great news! We have approved %s’s business banking application and your Wise account is ready to use.
Let's get started with your new account so you can begin to accept payments, get paid instantly, and enjoy the benefit
of no fees on most of your transactions. Your digital banking future is here!
<br /><br />
To get started, download the Wise mobile app from the <a href="https://wise.us/product/apps">https://wise.us/product/apps</a> and login to your account.
<br />
Note: You can view your account and routing number within the app.
<br /><br />
Your Wise business debit card will arrive in 7-10 business days. Please keep an eye out in the mail. The card will be
shipped to the mailing address you specified during signup:
<br /><br />
%s
<br /><br />
Need some help getting everything set up on Wise? We love helping business owners, so please reach out! You can head
to <a href="https://hub.wise.us">hub.wise.us</a>, where you will find product guides or chat to us if you have any
questions.
<br /><br />
Welcome to Wise!<br />
Wise Customer Success Team
<br /><br />
© 2020 Wise Company.<br />
Wise Company banking services provided by BBVA USA., member FDIC.<br />
The Wise Company Visa Card is issued pursuant to license from Visa U.S.A.<br />
<a href="https://hub.wise.us/en/collections/1642414-legal-disclosures-terms-and-policies"> Legal Disclosures, Terms and Policies </a><br />
Need help? Email us at chat@wise.us.
`

// BusinessDeclinedEmailSubject ..
const BusinessDeclinedEmailSubject = "Important information about your Wise Banking Application"

// BusinessDeclinedEmail business status declined copy
const BusinessDeclinedEmail = `Hi %s, Wise has carefully reviewed your application and unfortunately we are unable to approve %s's Wise Banking application at this time.<br />

Your application may have been rejected because we were unable to verify some of the information provided in the application or your business does not qualify for an account.<br />

We know you put a lot of time and effort into the application and this is not what you were hoping to read. We really did want to send you the other email...so if you would like to work with us to correct the issues preventing your account approval, please reach out to our Customer Success team at hub.wise.us or send an email to chat@wise.us and we’ll get in touch quickly!<br />
<br />
Thank you,<br />
Wise Customer Success Team<br />

© 2020 Wise Company.<br />
Wise Company banking services provided by BBVA USA.<br />
The Wise Company Visa Card is issued pursuant to license from Visa U.S.A. <br />
<a href="https://hub.wise.us/en/collections/1642414-legal-disclosures-terms-and-policies"> Legal Disclosures, Terms and Policies </a><br/>
Need help? Email us at chat@wise.us.
`

const PaymentReceivedSubject = `Transfer Receipt`

const PaymentReceivedEmail = `%s, <br /><br />
You have received a payment  <br /><br />
%s<br />
$%s<br />
%s<br /><br />
Funds will be available in your account within 24 hours.<br /><br />
Our customer success team can help you with your questions. Please visit hub.wise.us or send an email to chat@wise.us.<br /><br />
This is a service email from Wise. Please note that you may receive service emails in accordance with your Wise service agreements, whether or not you elect to receive promotional email. Please don't reply directly to this automatically generated email message.<br /><br />	Read our privacy policy at wise.us.<br /><br />
© 2020 Wise Company. Wise Company banking services provided by BBVA USA.<br /><br />`

const CheckPaymentReceivedEmail = `%s, <br /><br />
You were mailed a check.  <br /><br />
%s<br />
$%s<br />
%s<br /><br />
You should receive your check by mail in 3-5 business days.<br /><br />
Our customer success team can help you with your questions. Please visit hub.wise.us or send an email to chat@wise.us.<br /><br />
This is a service email from Wise. Please note that you may receive service emails in accordance with your Wise service agreements, whether or not you elect to receive promotional email. Please don't reply directly to this automatically generated email message.<br /><br />	Read our privacy policy at wise.us.<br /><br />
© 2020 Wise Company. Wise Company banking services provided by BBVA USA.<br /><br />`

const PaymentSentSubject = `Transfer Receipt`

const PaymentSentEmail = `%s, <br /><br />
You have successfully made a payment from your Wise account %s <br />
The funds have been withdrawn from your account and will appear in your transaction history.<br /><br />
%s<br />
%s<br />
$%s<br />
%s<br /><br />
Our customer success team can help you with your questions. Please visit hub.wise.us or send an email to chat@wise.us.<br /><br />
This is a service email from Wise. Please note that you may receive service emails in accordance with your Wise service agreements, whether or not you elect to receive promotional email. Please don't reply directly to this automatically generated email message.<br /><br />
Read our privacy policy at wise.us.<br /><br />
© 2020 Wise Company. Wise Company banking services provided by BBVA USA.<br /><br />`

const CustomerInvoiceSubject = `You have received an invoice sent by %s via Wise`
const CustomerResendInvoiceSubject = `REMINDER: You have received an invoice sent by %s via Wise`

const CustomerInvoiceEmail = `Hello %s,<br /><br />
You have received an invoice sent by %s via Wise<br /><br />
%s<br />
%s<br />
$%s<br />
<a href=%s>Pay %s $%s </a> <br /><br />
Once payment is complete, $%s will be available instantly in %s's Wise Banking account.<br /><br />
You're receiving this email because you have been invoiced by  %s, that banks with Wise, and Wise partners with Stripe to provide payments processing. Read our privacy policy at wise.us.<br />
© 2020 Wise Company. Wise Company banking services provided by BBVA USA.<br /><br />`

const CustomerInvoiceSMS = `Hello %s,
You have received an invoice sent by %s via Wise
%s
%s
$%s
To make the payment, click on the link below:
%s
Once payment is complete, $%s will be available instantly in %s's Wise Banking account.
You're receiving this email because you have been invoiced by  %s, that banks with Wise, and Wise partners with Stripe to provide payments processing. Read our privacy policy at wise.us.
© 2020 Wise Company. Wise Company banking services provided by BBVA USA.`

const CustomerInvoiceEmailWithoutPayURL = `Hello %s,<br /><br />
You have received an invoice sent by %s via Wise<br /><br />
%s<br />
%s<br />
$%s<br />
Once payment is complete, $%s will be available in %s's Wise Banking account.<br /><br />
You're receiving this email because you have been invoiced by  %s, that banks with Wise, and Wise partners with Stripe to provide payments processing. Read our privacy policy at wise.us.<br />
© 2020 Wise Company. Wise Company banking services provided by BBVA USA.<br /><br />`

const CustomerInvoiceSMSWithoutPayURL = `Hello %s,
You have received an invoice sent by %s via Wise
%s
%s
$%s
Once payment is complete, $%s will be available in %s's Wise Banking account.
You're receiving this email because you have been invoiced by  %s, that banks with Wise, and Wise partners with Stripe to provide payments processing. Read our privacy policy at wise.us.
© 2020 Wise Company. Wise Company banking services provided by BBVA USA.`

const BusinessInvoiceSubject = `You have sent an invoice of $%s to %s`

const BusinessInvoiceEmail = `Hello %s, <br /><br />
You have sent an invoice to %s via Wise.<br /><br />
%s<br />
%s<br />
$%s<br /><br />
Once payment is complete, $%s will be available instantly in %s's Wise Banking account.<br /><br />
This is a service email from Wise. Please note that you may receive service emails in accordance with your Wise service agreements, whether or not you elect to receive promotional email. Read our privacy policy at wise.us.<br />
© 2020 Wise Company. Wise Company banking services provided by BBVA USA.<br /><br />`

const CustomerReceiptSubject = `You have paid an invoice of $%s sent by %s`

const CustomerReceiptEmail = `Hello %s, <br /><br />
You have paid an invoice that was sent by %s via Wise. <br /><br />
%s<br />
%s<br />
$%s<br /><br />
Now that payment is complete, $%s is available in %s's Wise Banking account.<br /><br />
You're receiving this email because you have been invoiced by %s, that banks with Wise, and Wise partners with Stripe to provide payments processing. Read our privacy policy at wise.us.<br />
© 2020 Wise Company. Wise Company banking services provided by BBVA USA.<br /><br />`

const CustomerReceiptEmailWithInvoiceViewLink = `Hello %s, <br /><br />
You have paid an invoice that was sent by %s via Wise. <br /><br />
%s<br />
%s<br />
$%s<br /><br />
Here is <a href="%s" > the link </a> to invoice. <br />
Here is <a href="%s" > the link </a> to invoice receipt. <br />
Now that payment is complete, $%s is available in %s's Wise Banking account.<br /><br />
You're receiving this email because you have been invoiced by %s, that banks with Wise, and Wise partners with Stripe to provide payments processing. Read our privacy policy at wise.us.<br />
© 2020 Wise Company. Wise Company banking services provided by BBVA USA.<br /><br />`

const CustomerCardReaderReceiptSubject = `You've paid $%s to %s`

const CustomerCardReaderReceiptEmail = `Hello, <br /><br />
You've paid $%s to %s <br /><br />
%s<br />
%s<br />
$%s<br /><br />
Now that payment is complete, $%s is available in %s's Wise Banking account.<br /><br />
You're receiving this email because you have been invoiced by %s, that banks with Wise, and Wise partners with Stripe to provide payments processing. Read our privacy policy at wise.us.<br />
© 2020 Wise Company. Wise Company banking services provided by BBVA USA.<br /><br />`

const BusinessReceiptSubject = `Your invoice of $%s to %s has been paid!`

const BusinessReceiptEmail = `Hello %s, <br /><br />
Your invoice of $%s to %s has been paid! <br /><br />
%s<br />
%s<br />
$%s<br /><br />
Now that payment is complete, $%s is available in %s's Wise Banking account.<br /><br />
This is a service email from Wise. Please note that you may receive service emails in accordance with your Wise service agreements, whether or not you elect to receive promotional email. Read our privacy policy at wise.us.<br />
© 2020 Wise Company. Wise Company banking services provided by BBVA USA.<br /><br />`

const BusinessReceiptEmailWithInvoiceViewLink = `Hello %s, <br /><br />
Your invoice of $%s to %s has been paid! <br /><br />
%s<br />
%s<br />
$%s<br /><br />
Here is <a href="%s" > the link </a> to invoice. <br />
Now that payment is complete, $%s is available in %s's Wise Banking account.<br /><br />
This is a service email from Wise. Please note that you may receive service emails in accordance with your Wise service agreements, whether or not you elect to receive promotional email. Read our privacy policy at wise.us.<br />
© 2020 Wise Company. Wise Company banking services provided by BBVA USA.<br /><br />`

const BusinessCardReaderReceiptSubject = `You've been paid $%s via card reader!`

const BusinessCardReaderReceiptEmail = `Hello %s, <br /><br />
You've been paid $%s via card reader! <br /><br />
%s<br />
%s<br />
$%s<br /><br />
Now that payment is complete, $%s is available in %s's Wise Banking account.<br /><br />
This is a service email from Wise. Please note that you may receive service emails in accordance with your Wise service agreements, whether or not you elect to receive promotional email. Read our privacy policy at wise.us.<br />
© 2020 Wise Company. Wise Company banking services provided by BBVA USA.<br /><br />`

const TransferRequestSubject = `Incoming Payment from %s`

const TransferRequestEmail = `%s,<br/><br/>
%s uses Wise’s instant business payment service and would like to pay you. <br/><br/>
To accept the payment, click on the link below: %s<br/><br/>
Our customer success team can help you with your questions. Please visit hub.wise.us or send an email to chat@wise.us.<br /><br />
This is a service email from Wise. Please note that you may receive service emails in accordance with your Wise service agreements, whether or not you elect to receive promotional email. Please don't reply directly to this automatically generated email message.<br /><br />
Read our privacy policy at wise.us.<br /><br />
© 2020 Wise Company. Wise Company banking services provided by BBVA USA.<br /><br />`
