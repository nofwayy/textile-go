package core

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (a *api) lsContacts(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	var contacts []ContactInfo

	threadId := opts["thread"]
	if threadId == "default" {
		threadId = a.node.config.Threads.Defaults.ID
	}
	if threadId != "" {
		thrd := a.node.Thread(threadId)
		if thrd == nil {
			g.String(http.StatusNotFound, ErrThreadNotFound.Error())
			return
		}
		contacts = make([]ContactInfo, 0)
		for _, p := range thrd.Peers() {
			contact := a.node.Contact(p.Id)
			if contact != nil {
				contacts = append(contacts, *contact)
			}
		}
	} else {
		contacts, err = a.node.Contacts()
	}

	g.JSON(http.StatusOK, contacts)
}

func (a *api) getContacts(g *gin.Context) {
	id := g.Param("id")

	info := a.node.Contact(id)
	if info == nil {
		g.String(http.StatusNotFound, "contact not found")
		return
	}

	g.JSON(http.StatusOK, info)
}

func (a *api) addContacts(g *gin.Context) {
	args, err := a.readArgs(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	if len(args) < 2 {
		g.String(http.StatusBadRequest, "missing peer id or address")
		return
	}
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	id := args[0]
	if err := a.node.AddContact(id, args[1], opts["username"]); err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	info := a.node.Contact(id)
	if info == nil {
		g.String(http.StatusNotFound, "contact not created")
		return
	}

	g.JSON(http.StatusCreated, info)
}

func (a *api) searchContacts(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	local, err := strconv.ParseBool(opts["local"])
	if err != nil {
		local = false
	}
	lucky, err := strconv.ParseBool(opts["lucky"])
	if err != nil {
		lucky = false
	}
	limit, err := strconv.Atoi(opts["limit"])
	if err != nil {
		limit = 5
	}
	wait, err := strconv.Atoi(opts["wait"])
	if err != nil {
		wait = 5
	}
	query := &ContactInfoQuery{
		Id:       opts["peer"],
		Address:  opts["address"],
		Username: opts["username"],
		Local:    local,
		Lucky:    lucky,
		Limit:    limit,
		Wait:     wait,
	}

	infos, err := a.node.FindContact(query)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.JSON(http.StatusOK, infos)
}
