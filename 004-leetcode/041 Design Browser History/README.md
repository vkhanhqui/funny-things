[1472. Design Browser History](https://leetcode.com/problems/design-browser-history)

```cs
public class BrowserHistory
{
    public int currentPage = 0;
    public int real_steps = 1;
    public List<string> homepages = new List<string>();

    public BrowserHistory(string homepage)
    {
        homepages.Add(homepage);
    }

    public void Visit(string url)
    {
        if(currentPage + 2 > homepages.Count){
            homepages.Add(url);
        }
        else {
            homepages[currentPage+1] = url;
        }
        currentPage += 1;
        real_steps = currentPage + 1;
    }

    public string Back(int steps)
    {
        currentPage = Math.Max(currentPage - steps, 0);
        return homepages[currentPage];
    }

    public string Forward(int steps)
    {
        currentPage = Math.Min(real_steps - 1, currentPage + steps);
        return homepages[currentPage];
    }
}

```