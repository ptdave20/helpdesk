package main

import (
	"encoding/json"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type (
	Config struct {
		RedirectURI  string `json:"redirect_uri"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}
	SimpleResult struct {
		Result bool `json:"result"`
	}
	Device struct {
		Id          bson.ObjectId `bson:"_id"`
		BuildingId  bson.ObjectId `bson:"building_id,omitempty"`
		DomainId    bson.ObjectId `bson:"domain_id,omitempty"`
		LastUpdate  time.Time     `bson:"last_update,omitempty"`
		FirstUpdate time.Time     `bson:"first_update,omitempty"`
	}
	DeviceEvent struct {
		When  time.Time     `bson:"when"`
		What  string        `bson:"what"`
		WhoId bson.ObjectId `bson:"who_id"`
	}
	Group struct {
		Id      bson.ObjectId   `bson:"_id"`
		Members []bson.ObjectId `bson:"members"`
	}
	Role struct {
		DomainAdmin     bool `bson:"domain_admin"`
		DomainSetRole   bool `bson:"domain_set_role"`
		DomainModDep    bool `bson:"domain_mod_dep"`
		DomainModBldg   bool `bson:"domain_mod_bldg"`
		DomainModUser   bool `bson:"domain_mod_user"`
		DomainModConfig bool `bson:"domain_mod_config"`
		BldgViewTicket  bool `bson:"bldg_view_tickets"`
	}
	Building struct {
		Id   bson.ObjectId `bson:"_id"`
		Name string        `bson:"name"`
	}
	EmailTemplate struct {
		SubmitTicket   string `bson:"submit_ticket"`   // When a user creates a ticket, what the submitter sees
		NewTicket      string `bson:"new_ticket"`      // When a user creates a ticket, what department sees
		AssignedTicket string `bson:"assigned_ticket"` // When a ticket has been assigned, what the user sees
		SolvedTicket   string `bson:"solved_ticket"`   // When a ticket has been solved, what the submitter sees
		NotedTicket    string `bson:"noted_ticket"`    // When a ticket has a public note, what the submitter sees
		DocumentAdded  string `bson:"document_added"`  // When a ticket has a new document added, universal
	}
	DomainDefault struct {
		UserRole Role `bson:"user_roles"`
	}
	DomainSettings struct {
		Defaults       DomainDefault `bson:"domain_defaults"`
		EmailTemplates EmailTemplate `bson:"email_templates"`
		KeepStats      bool          `bson:"keep_user_stats"`
	}
	Domain struct {
		Id              bson.ObjectId  `bson:"_id"`
		Name            string         `bson:"name"`
		AcceptedDomains []string       `bson:"accepted_domains"`
		Settings        DomainSettings `bson:"domain_settings"`
		Buildings       []Building     `bson:"buildings,omitempty"`
		Departments     []Department   `bson:"departments,omitempty"`
	}
	Note struct {
		Id        bson.ObjectId `bson:"_id"`
		Public    bool          `bson:"public"`
		Submitter bson.ObjectId `bson:"submitter"`
		Created   time.Time     `bson:"created"`
		Detail    string        `bson:"details"`
	}
	Category struct {
		Id   bson.ObjectId `bson:"_id"`
		Name string        `bson:"name"`
	}
	DepartmentUser struct {
		UserId          bson.ObjectId `bson:"user_id"`
		DepAdmin        bool          `bson:"dep_admin"`
		DepAssignTicket bool          `bson:"dep_assign_ticket"`
		DepCloseTicket  bool          `bson:"dep_close_ticket"`
	}
	Department struct {
		Id                 bson.ObjectId    `bson:"_id"`
		Name               string           `bson:"name"`
		Category           []Category       `bson:"category,omitempty"`
		IsBuildingSpecific bool             `bson:"is_building_specific"`
		Building           bson.ObjectId    `bson:"building,omitempty"`
		Users              []DepartmentUser `bson:"department_users,omitempty"`
	}
	DepartmentUpdate struct {
		Name               string        `json:"Name"`
		IsBuildingSpecific bool          `json:"IsBuildingSpecific"`
		Building           bson.ObjectId `json:"Building"`
	}
	Document struct {
		Id        bson.ObjectId `bson:"_id,omitempty"`
		Created   time.Time     `bson:"created"`
		Submitter bson.ObjectId `bson:"submitter,omitempty"`
		Name      string        `bson:"name"`
		Data      []byte        `bson:"data"`
		Mime      string        `bson:"mime"`
	}
	TicketStatus struct {
		Value string `bson:"value"`
		Name  string `bson:"name"`
	}
	TicketMessage struct {
		From    bson.ObjectId `bson:"from,omitempty",json:"omitempty"`
		When    time.Time     `bson:"when",json:"omitempty"`
		Message string        `bson:"message"`
	}
	Ticket struct {
		Id          bson.ObjectId   `bson:"_id,omitempty"`
		Domain      bson.ObjectId   `bson:"domain_id"`
		Submitter   bson.ObjectId   `bson:"submitter,omitempty"`
		AssignedTo  bson.ObjectId   `bson:"assigned_to,omitempty"`
		AttachedTo  bson.ObjectId   `bson:"attached_to,omitempty"`
		Building    bson.ObjectId   `bson:"building,omitempty"`
		Department  bson.ObjectId   `bson:"department,omitempty"`
		Category    bson.ObjectId   `bson:"category,omitempty"`
		Target      bson.ObjectId   `bson:"target,omitempty"`
		Subject     string          `bson:"subject,omitempty"`
		Created     time.Time       `bson:"created,omitempty"`
		Closed      time.Time       `bson:"closed,omitempty"`
		Updated     time.Time       `bson:"updated,omitempty"`
		Status      string          `bson:"status,omitempty"`
		Duration    time.Duration   `bson:"duration,omitempty"`
		Notes       []Note          `bson:"notes,omitempty"`
		Documents   []Document      `bson:"documents,omitempty"`
		Description string          `bson:"description,omitempty"`
		Solution    string          `bson:"solution,omitempty"`
		Messages    []TicketMessage `bson:"messages,omitempty"`
	}
	TicketUpdate struct {
		Department  bson.ObjectId `json:"Department,omitempty"`
		Category    bson.ObjectId `json:"Category,omitempty"`
		Description string        `json:"Description,omitempty"`
		Subject     string        `json:"Subject,omitempty"`
	}
	TicketCount struct {
		Day        int `bson:"day"`
		Month      int `bson:"month"`
		Year       int `bson:"year"`
		Submitted  int `bson:"submitted"`
		Closed     int `bson:"closed"`
		Assigned   int `bson:"assigned"`
		AssignedTo int `bson:"assigned_to"`
		Noted      int `bson:"noted"`
	}
	Session struct {
		Id     bson.ObjectId `bson:"_id,omitempty"`
		UserId bson.ObjectId `bson:"user_id,omitempty"`
	}
	User struct {
		Id         bson.ObjectId   `bson:"_id,omitempty",json:"omitempty"`
		Domain     bson.ObjectId   `bson:"domain_id,omitempty",json:"omitempty"`
		Department []bson.ObjectId `bson:"department,omitempty",json:"omitempty"`
		Building   bson.ObjectId   `bson:"location,omitempty",json:"omitempty"`

		NewUser bool `bson:"new_user",json:"omitempty"`

		GoogleId  string `bson:"google_id",json:"omitempty"`
		Firstname string `bson:"firstname",json:"omitempty"`
		Lastname  string `bson:"lastname",json:"omitempty"`
		Email     string `bson:"email",json:"omitempty"`
		Picture   string `bson:"picture",json:"omitempty"`
		Room      string `bson:"room",json:"omitempty"`
		Phone     string `bson:"phone",json:"omitempty"`
		Extension string `bson:"extension",json:"omitempty"`

		Enabled bool `bson:"enabled",json:"omitempty"`

		Roles Role `bson:"role",json:"omitempty"`

		FirstLogin time.Time `bson:"first_login",json:"omitempty"`
		LastLogin  time.Time `bson:"last_login",json:"omitempty"`
		RolesSet   time.Time `bson:"role_set",json:"omitempty"`

		TicketStats []TicketCount   `bson:"ticket_count",json:"omitempty"`
		Submitted   []bson.ObjectId `bson:"submitted",json:"omitempty"`
	}
	UserUpdate struct {
		Firstname string
		Lastname  string
		Email     string
		Room      string
		Phone     string
		Extension string
		Building  bson.ObjectId
	}
	GoogleUserV2 struct {
		Id            string `json:"id,omitempty"`
		Email         string `json:"email,omitempty"`
		VerifiedEmail bool   `json:"verified_email,omitempty"`
		GivenName     string `json:"given_name,omitempty"`
		FamilyName    string `json:"family_name,omitempty"`
		Picture       string `json:"picture,omitempty"`
		Hd            string `json:"hd,omitempty"`
	}
	MassMessage struct {
		Id      string        `bson:"_id"`
		Domain  bson.ObjectId `bson:"domain"`
		Expired bool          `bson:"expired"`
		Message string        `bson:"message"`
		From    bson.ObjectId `bson:"from"`
	}
	News struct {
		Id      string        `bson:"_id"`
		Domain  bson.ObjectId `bson:"domain"`
		Show    bool          `bson:"show"`
		From    bson.ObjectId `bson:"from"`
		Message string        `bson:"message"`
	}
)

func (t Ticket) Marshal() string {
	b, _ := json.Marshal(t)
	return string(b)
}

func (d Department) GetMember(id bson.ObjectId) (*DepartmentUser, error) {
	for i := 0; i < len(d.Users); i++ {
		if d.Users[i].UserId == id {
			return &d.Users[i], nil
		}
	}
	return nil, fmt.Errorf("Department: user %s is not a member of \"%s\"", id.Hex(), d.Name)
}

func (d Department) CanEditTicket(u User, t Ticket) bool {
	var depUser int = -1

	// Is the user in the department on their list?
	if !u.InDepartment(t.Department) {
		return false
	}

	// Find the user in the department to get their roles
	for i := 0; i < len(d.Users); i++ {
		if d.Users[i].UserId.Hex() == u.Id.Hex() {
			depUser = i
			break
		}
	}

	// Can't find them? then they can only view... maybe
	if depUser == -1 {
		return false
	}

	//
	return true
}

func (u User) Marshal() ([]byte, error) {
	ret, err := json.Marshal(u)
	return ret, err
}
func (u User) GetAssignedTickets(db *mgo.Database) ([]Ticket, error) {
	var t []Ticket
	c := db.C(TicketsC)
	err := c.Find(bson.M{"assigned_to": u.Id}).All(&t)
	return t, err
}
func (u User) InDepartment(id bson.ObjectId) bool {
	// If it is not a valid object id then return false
	if !id.Valid() {
		return false
	}

	// To limit the cpu work when retrieving a hex
	hex := id.Hex()

	// Iterate through the departments that the user is associated with
	for i := 0; i < len(u.Department); i++ {
		if hex == u.Department[i].Hex() {
			return true
		}
	}
	return false
}
func (u User) CanViewTicket(ticket Ticket) bool {
	if ticket.Submitter.Hex() == u.Id.Hex() {
		return true
	}

	if ticket.AssignedTo.Hex() == u.Id.Hex() {
		return true
	}

	if u.InDepartment(ticket.Department) {
		return true
	}

	if u.Roles.DomainAdmin {
		return true
	}

	if u.Roles.BldgViewTicket && ticket.Building.Hex() == u.Building.Hex() {
		return true
	}

	return false
}
func (u User) CanEditTicket(ticket Ticket) bool {
	if ticket.AssignedTo.Hex() == u.Id.Hex() {
		return true
	}

	if ticket.Submitter.Hex() == u.Id.Hex() {
		return true
	}

	if u.Roles.DomainAdmin {
		return true
	}

	return false
}
func (u User) CanAddNote(ticket Ticket) bool {
	if ticket.AssignedTo.Hex() == u.Id.Hex() {
		return true
	}
	if ticket.Submitter.Hex() == u.Id.Hex() {
		return true
	}
	return u.InDepartment(ticket.Department)
}
func (u User) CanAddMessage(ticket Ticket) bool {
	if ticket.AssignedTo.Hex() == u.Id.Hex() {
		return true
	}
	if ticket.Submitter.Hex() == u.Id.Hex() {
		return true
	}
	if u.Roles.DomainAdmin {
		return true
	}
	return u.InDepartment(ticket.Department)
}
func (u User) CanAssign(dep Department, ticket Ticket) bool {
	if u.Roles.DomainAdmin {
		return true
	}
	for i := 0; i < len(dep.Users); i++ {
		if dep.Users[i].UserId.Hex() == u.Id.Hex() && dep.Users[i].DepAssignTicket {
			return true
		}
	}
	return false
}
func (u User) CanDelete(ticket Ticket) bool {
	return false
}
func (u User) CanUpdate(ticket Ticket) bool {
	return false
}
func (u User) FindUserById(id string, db *mgo.Database) bool {
	c := db.C(UsersC)
	err := c.Find(bson.M{"_id": id}).One(&u)
	if err != nil {
		return false
	}
	return true
}
func (r SimpleResult) Marshal() string {
	b, _ := json.Marshal(r)
	return string(b)
}
