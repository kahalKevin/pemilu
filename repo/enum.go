package repo

// ==adm1nb0sque__

type Tingkat int

const (
   RI    Tingkat = 0
   DPD    Tingkat = 1
   DPRRI   Tingkat = 2
   DPRD1 Tingkat = 3
   DPRD2  Tingkat = 4
)

func (tingkat Tingkat) String() string {
    names := [...]string{
        "RI", 
        "DPD", 
        "DPRRI", 
        "DPRD1",
        "DPRD2"}
    if tingkat < RI || tingkat > DPRD2 {
      return "Unknown"
    }
    return names[tingkat]
}


type Role int

const (
   ADMIN    Role = 0
   CALON    Role = 1
)

func (role Role) String() string {
    names := [...]string{
        "ADMIN", 
        "CALON"}
    if role < ADMIN || role > CALON {
      return "Unknown"
    }
    return names[role]
}