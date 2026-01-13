import JINN1

try:
    JINN1.main() 
except AttributeError:
    try:
        JINN1.menu()
    except:
        pass
      
