import React from 'react'
import ProductsContainer from './ProductsContainer'
import CartContainer from './CartContainer'

import { Tab } from 'semantic-ui-react'



const panes = [
        { menuItem: 'Products', render: () => <Tab.Pane><ProductsContainer /></Tab.Pane> },
        { menuItem: 'Card', render: () => <Tab.Pane><CartContainer /></Tab.Pane> }
 ]


const App = () => (
  <div>
    <h2>Shopping Cart Example</h2>
    <hr/>
    <Tab panes={panes} />
  </div>
)

export default App
