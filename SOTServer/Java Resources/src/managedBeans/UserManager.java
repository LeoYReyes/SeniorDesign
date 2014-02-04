package managedBeans;
import javax.ejb.EJB;
import javax.faces.bean.ManagedBean;
import javax.faces.bean.SessionScoped;
import javax.xml.registry.infomodel.User;

@SessionScoped
@ManagedBean(name="user")
public class UserManager {
	private String userName;
	private String password;
	private User current;
	
	//@EJB
	//private UserService userService;
}
