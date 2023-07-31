INSERT INTO public.ms_payment_method (id,"name",description,code) VALUES
	 (25,'BCA Pak Paul','BCA Pak Paul','bca_paul'),
	 (26,'BCA Pak Pras','BCA Pak Pras','bca_pras'),
	 (27,'BCA Pak Mike','BCA Pak Mike','bca_mike');

INSERT INTO public.ms_payment_status (xms_id,"name",code) VALUES
	 (7,'Pending','PENDING'),
	 (4,'Paid','PAID'),
	 (5,'Partially Paid','PARTIALLY_PAID'),
	 (9,'Unpaid','UNPAID'),
	 (2,'Expired','EXPIRED');

INSERT INTO public.ms_order_status (id,"name") VALUES
	 (1,'On Process'),
	 (2,'Completed'),
	 (3,'Cancelled'),
	 (4,'Draft'),
	 (5,'Pending'),
	 (6,'Verification');

INSERT INTO public.ms_order_integration_status (id,"name") VALUES
	 (1,'Unfulfilled'),
	 (2,'Fulfilled'),
	 (3,'Rejected'),
	 (4,'Cancelled'),
	 (5,'Follow Up Consign');

INSERT INTO public.ms_fulfillment_status (id,"name") VALUES
	 (1,'Waiting'),
	 (2,'Ready for Picking'),
	 (3,'Ready for Packing'),
	 (5,'Request Shipping'),
	 (6,'Ready for Shipping'),
	 (7,'Rejected'),
	 (8,'Complete'),
	 (9,'Hold'),
	 (4,'Ready for PDI');

INSERT INTO public.ms_office (id,"name",code,"type",address,province_id,district_id,sub_district_id,zip,status) VALUES
	 (1,'Store Jakarta','jkt','offline_store','Jalan Sultan Iskandar Muda No 18',6,56,711,12440,1),
	 (2,'Store Surabaya','sby','offline_store','Bukit Darmo Golf Boulevard Office Park I No 9-10',11,175,2777,60226,1),
	 (3,'Warehouse Jakarta','wh-jkt','warehouse','Jalan Sultan Iskandar Muda No 18',6,56,711,12440,1),
	 (4,'Warehouse Surabaya','wh-sby','warehouse','Bukit Darmo Golf Boulevard Office Park I No 9-10',11,175,2777,60226,1),
	 (5,'Meru','meru','consign','Meru Luxury, Jalan tukas barito no 23, Panjer Denpasar, Bali',2,27,310,80225,1),
	 (6,'HGS','hgs','consign','Handi Ruko Blok K1, Mall of Indonesia',6,60,745,14240,1),
	 (7,'Anya','anya','consign','Katamaran indah 3 no 2 E ',6,60,745,14460,1),
	 (8,'Angela','','consign','Muara Karang D 10 selatan No 16',6,60,748,14450,1),
	 (9,'Kick Ave','','consign','Jl. Deperdag 1 blok H no 10, RT.7/RW.4, Gandaria Utara',6,58,726,12140,1),
	 (10,'Marisa T','','consign','1 Park Residences tower B unit 6F, Jl. Kyai Moh. Syafii Hadzami no 1, RT 1 / RW 10',6,58,724,12440,1),
	 (11,'Winson','','consign','Beast Kicks Store Jl Rusa No 9-11',10,131,2028,50123,1),
	 (12,'Eunike','','consign','Jl. Dr. Cipto no 119 (E wholesale)',6,56,711,12440,1),
	 (13,'Inggrid (Snowceline)','','consign','Gading Griya Lestari blok a6 no 2',6,60,744,14140,1),
	 (14,'Jamtangan','','consign','Jamtangan',6,56,711,12440,1),
	 (16,'Bazaar JKT','bazaar-jkt','offline_store','Jakarta',6,56,711,12440,1);

INSERT INTO public.ms_courier (id,group_id,short_name,long_name,group_name,service) VALUES
	 (1,1,'JNE','','JNE','Yes'),
	 (2,1,'JNE','','JNE','Reg'),
	 (3,1,'JNE','','JNE','Yes Cashless'),
	 (4,1,'JNE','','JNE','OKE Cashless'),
	 (5,1,'JNE','','JNE','Reg Cashless'),
	 (6,1,'JNE','','JNE','JOB'),
	 (7,1,'JNE','','JNE','JTR (JNE Trucking)'),
	 (8,2,'J&T','','J&T','Reg'),
	 (9,2,'J&T','','J&T','JND JSD'),
	 (10,3,'Si Cepat','','SICEPAT','Reg'),
	 (11,4,'TIKI','','TIKI',''),
	 (12,4,'TIKI','','TIKI','ONS'),
	 (13,5,'Anteraja','','ANTERAJA','Reg'),
	 (14,6,'Paxel','','PAXEL',''),
	 (15,7,'Go-Send','','GOSEND','Instant'),
	 (16,7,'Go-Send','','GOSEND','Sameday'),
	 (17,8,'Grab-Express','','GRABEXPRESS','Instant'),
	 (18,8,'Grab-Express','','GRABEXPRESS','Sameday'),
	 (19,9,'Self Pick Up','','SELF_PICK_UP',''),
	 (20,10,'Kurir Pribadi','','INTERNAL_COURIER',''),
	 (21,11,'Kurir Lainnya','','OTHERS',''),
	 (22,2,'J&T','','J&T','Reg Cashless'),
	 (23,3,'Si Cepat','','SICEPAT','Reg Cashless'),
	 (24,3,'Si Cepat','','SICEPAT','Halu Cashless'),
	 (25,2,'J&T','','J&T','Economy Cashless'),
	 (26,1,'JNE','','JNE','Superspeed');
