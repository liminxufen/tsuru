package service_test

import (
	. "github.com/timeredbull/tsuru/api/app"
	. "github.com/timeredbull/tsuru/api/service"
	. "launchpad.net/gocheck"
	. "github.com/timeredbull/tsuru/database"
)

func (s *ServiceSuite) createService() {
	s.serviceType = &ServiceType{Name: "Mysql", Charm: "mysql"}
	s.serviceType.Create()

	s.service = &Service{ServiceTypeId: s.serviceType.Id, Name: "my_service"}
	s.service.Create()
}

func (s *ServiceSuite) TestGetService(c *C) {
	s.createService()
	id := s.service.Id
	sTypeId := s.service.ServiceTypeId
	s.service.Id = ""
	s.service.ServiceTypeId = ""
	s.service.Get()

	c.Assert(s.service.Id, Equals, id)
	c.Assert(s.service.ServiceTypeId, Equals, sTypeId)
}

func (s *ServiceSuite) TestAllServices(c *C) {
	st := ServiceType{Name: "mysql", Charm: "mysql"}
	st.Create()
	se := Service{ServiceTypeId: st.Id, Name: "myService"}
	se2 := Service{ServiceTypeId: st.Id, Name: "myOtherService"}
	se.Create()
	se2.Create()

	s_ := Service{}
	results := s_.All()
	c.Assert(len(results), Equals, 2)
}

func (s *ServiceSuite) TestCreateService(c *C) {
	s.createService()
	se := Service{Id: s.service.Id}
	se.Get()

	c.Assert(se.Id, Equals, s.service.Id)
	c.Assert(se.ServiceTypeId, Equals, s.serviceType.Id)
	c.Assert(se.Name, Equals, s.service.Name)
}

func (s *ServiceSuite) TestDeleteService(c *C) {
	s.createService()
	s.service.Delete()

	rows, err := Db.Query("SELECT count(*) FROM services WHERE name = 'my_service'")
	c.Assert(err, IsNil)

	var qtd int
	for rows.Next() {
		rows.Scan(&qtd)
	}

	c.Assert(qtd, Equals, 0)
}

func (s *ServiceSuite) TestRetrieveAssociateServiceType(c *C) {
	serviceType := ServiceType{Name: "Mysql", Charm: "mysql"}
	serviceType.Create()

	service := &Service{
		ServiceTypeId: serviceType.Id,
		Name:          "my_service",
	}
	service.Create()
	retrievedServiceType := service.ServiceType()

	c.Assert(serviceType.Id, Equals, retrievedServiceType.Id)
	c.Assert(serviceType.Name, Equals, retrievedServiceType.Name)
	c.Assert(serviceType.Charm, Equals, retrievedServiceType.Charm)
}

func (s *ServiceSuite) TestBindService(c *C) {
	s.createService()
	app := &App{Name: "my_app", Framework: "django"}
	app.Create()
	s.service.Bind(app)

	rows, err := Db.Query("SELECT service_id, app_id FROM service_apps WHERE service_id = ? AND app_id = ?", s.service.Id, app.Id)
	c.Assert(err, IsNil)

	var serviceId int64
	var appId int64
	for rows.Next() {
		rows.Scan(&serviceId, &appId)
	}

	c.Assert(s.service.Id, Equals, serviceId)
	c.Assert(app.Id, Equals, appId)
}

func (s *ServiceSuite) TestUnbindService(c *C) {
	serviceType := &ServiceType{Name: "Mysql", Charm: "mysql"}
	serviceType.Create()

	service := &Service{ServiceTypeId: s.serviceType.Id, Name: "my_service"}
	service.Create()
	app := &App{Name: "my_app", Framework: "django"}
	app.Create()
	service.Bind(app)
	service.Unbind(app)

	query := make(map[string]interface{})
	query["service_id"] = service.Id
	query["app_id"] = app.Id

	collection := Mdb.C("service_apps")
	qtd, err := collection.Find(query).Count()
	// rows, err := Db.Query("SELECT count(*) FROM service_apps WHERE service_id = ? AND app_id = ?", s.service.Id, app.Id)
	c.Assert(err, IsNil)

	// var qtd int
	// for rows.Next() {
	// 	rows.Scan(&qtd)
	// }

	c.Assert(qtd, Equals, 0)
}
